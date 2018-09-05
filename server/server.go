package main

import (
	"encoding/json"
	"flag"
	"fmt"
	influxdb "github.com/influxdata/influxdb/client/v2"
	pb "github.com/mnbbrown/healthsignal/healthsignal"
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
)

type healthSignalServer struct {
	influxClient influxdb.Client
	pingChan     chan *pb.Ping
}

func preparePings(batch influxdb.BatchPoints, pings []*pb.Ping) {
	for _, ping := range pings {
		tags := map[string]string{
			"location": ping.Location,
			"endpoint": strconv.FormatInt(int64(ping.Endpoint), 10),
		}
		fields := map[string]interface{}{
			"response": ping.RequestDuration,
			"system":   ping.ConnectionDuration,
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
			log.Println("Here")
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
	qs := fmt.Sprintf("SELECT mean(\"response\") FROM ping WHERE (\"endpoint\"='%d') AND time >= now() - 12h GROUP BY time(10s) fill(0)", endpoint)
	log.Println(qs)
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
	Timestamp time.Time `json:"timestamp"`
	Value     int64     `json:"value"`
}

func (h *healthSignalServer) query(rw http.ResponseWriter, req *http.Request) {
	endpoint := req.URL.Query().Get("endpoint")
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
	var points []point
	for _, result := range response {
		for _, series := range result.Series {
			for _, raw := range series.Values {
				log.Println(raw)
				if timestamp, ok := raw[0].(json.Number); ok {
					if timeMs, err := timestamp.Int64(); err == nil {
						if rawValue, ok := raw[1].(json.Number); ok {
							if value, err := rawValue.Int64(); err == nil {
								ts := time.Unix(0, timeMs*int64(time.Millisecond))
								points = append(points, point{
									Timestamp: ts,
									Value:     value,
								})
							}
						}
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
	http.HandleFunc("/", h.query)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *webport), nil))
}

func newHealthSignalServer() (*healthSignalServer, error) {
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
	go grpcServer.Serve(lis)
	server.startWeb()
}
