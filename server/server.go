package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	pb "github.com/mnbbrown/healthsignal/healthsignal"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	port    = flag.Int("port", 10000, "HealthSignal API Port")
	webport = flag.Int("webport", 8080, "HealthSignal Web API Port")
	dbDsn   = flag.String("db", "user=healthsignal password=healthsignal dbname=healthsignal sslmode=disable", "Connection string for postgres instance")
)

type healthSignalServer struct {
	db           *sqlx.DB
	influxClient influxdb.Client
	pingChan     chan *pb.Ping
}

func preparePings(batch influxdb.BatchPoints, pings []*pb.Ping) {
	for _, ping := range pings {
		tags := map[string]string{
			"location": ping.Location,
			"endpoint": strconv.FormatInt(int64(ping.Endpoint), 10),
			"timedOut": strconv.FormatBool(ping.TimedOut),
		}
		fields := map[string]interface{}{
			"requestTime":    ping.RequestDuration,
			"connectionTime": ping.ConnectionDuration,
			"status":         ping.OnlineStatus,
		}
		pt, _ := influxdb.NewPoint("ping", tags, fields, time.Now())
		batch.AddPoint(pt)
	}
}

func (h *healthSignalServer) startSink() {
	ticker := time.NewTicker(time.Second * 5)
	pings := []*pb.Ping{}
	for {
		select {
		case ping := <-h.pingChan:
			pings = append(pings, ping)
		case <-ticker.C:
			if len(pings) == 0 {
				continue
			}
			log.Printf("Sending %d points", len(pings))
			bp, _ := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
				Database:  "pings",
				Precision: "s",
			})
			preparePings(bp, pings)
			pings = []*pb.Ping{}
			err := h.influxClient.Write(bp)
			if err != nil {
				log.Printf("Failed to write to influx: %s", err)
			}
		}
	}
}

func (h *healthSignalServer) getpoints(endpoint int) (res []influxdb.Result, err error) {
	qs := fmt.Sprintf("SELECT mean(\"requestTime\"), mean(\"connectionTime\"), last(\"status\") FROM ping WHERE (\"endpoint\"='%d') AND time >= now() - 6h GROUP BY time(30s),location fill(none)", endpoint)
	q := influxdb.NewQuery(qs, "pings", "ms")
	if response, err := h.influxClient.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

type point struct {
	Timestamp      int64       `json:"timestamp"`
	ConnectionTime json.Number `json:"connectionTime"`
	RequestTime    json.Number `json:"requestTime"`
	Status         json.Number `json:"status"`
	Location       string      `json:"location"`
}

func (h *healthSignalServer) query(rw http.ResponseWriter, req *http.Request) {
	endpoint := chi.URLParam(req, "endpointID")
	if endpoint == "" {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	endpointID, err := strconv.Atoi(endpoint)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	response, err := h.getpoints(endpointID)
	if err != nil {
		log.Println(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	points := make([]point, 0)
	for _, result := range response {
		for _, series := range result.Series {
			for _, raw := range series.Values {
				if timestamp, ok := raw[0].(json.Number); ok {
					if timeMs, err := timestamp.Int64(); err == nil {
						ts := time.Unix(0, timeMs*int64(time.Millisecond))
						point := point{
							Timestamp:      ts.Unix(),
							ConnectionTime: raw[1].(json.Number),
							RequestTime:    raw[2].(json.Number),
							Status:         raw[3].(json.Number),
						}
						if location, ok := series.Tags["location"]; ok {
							point.Location = location
						}
						points = append(points, point)
					}
				}
			}
		}
	}
	rw.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(points)
	if err != nil {
		log.Println(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rw.Write(b)
}

func (h *healthSignalServer) startWeb() {
	r := chi.NewRouter()
	r.Get("/endpoints", func(rw http.ResponseWriter, req *http.Request) {
		endpoints, err := h.listEndpoints()
		if err != nil {
			log.Println(err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(endpoints)
		if err != nil {
			log.Println(err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(j)
	})
	r.Get("/endpoints/{endpointID}/data", h.query)
	log.Printf("HTTP API listening on %d", *webport)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *webport), cors.Default().Handler(r)))
}

func newHealthSignalServer() (*healthSignalServer, error) {
	db, err := sqlx.Connect("postgres", *dbDsn)
	if err != nil {
		return nil, err
	}
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "healthsignal",
		Password: "healthsignal",
	})
	if err != nil {
		return nil, err
	}
	_, version, err := c.Ping(time.Second * 30)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to influxdb %s", version)
	server := &healthSignalServer{
		db:           db,
		influxClient: c,
		pingChan:     make(chan *pb.Ping),
	}
	go server.startSink()
	return server, nil
}

func (h *healthSignalServer) SavePing(ctx context.Context, ping *pb.Ping) (*pb.Empty, error) {
	log.Println("New Ping")
	h.pingChan <- ping
	return &pb.Empty{}, nil
}

func (h *healthSignalServer) listEndpoints() ([]*pb.Endpoint, error) {
	rows, err := h.db.Query("SELECT id, url, expectedstatus, name FROM endpoints")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var endpoints []*pb.Endpoint
	for rows.Next() {
		endpoint := &pb.Endpoint{}
		err = rows.Scan(&endpoint.Id, &endpoint.Url, &endpoint.ExpectedStatus, &endpoint.Name)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}
	return endpoints, nil
}

func (h *healthSignalServer) GetEndpoints(ctx context.Context, query *pb.EndpointsQuery) (*pb.Endpoints, error) {
	endpoints, err := h.listEndpoints()
	log.Printf("Found %d endpoints", len(endpoints))
	return &pb.Endpoints{Endpoints: endpoints}, err
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	server, err := newHealthSignalServer()
	if err != nil {
		panic(err)
	}
	pb.RegisterHealthSignalServer(grpcServer, server)
	log.Printf("gRPC listening on %d", *port)
	go grpcServer.Serve(lis)
	server.startWeb()
}
