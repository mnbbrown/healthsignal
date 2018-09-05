package main

import (
	"flag"
	"fmt"
	influxdb "github.com/influxdata/influxdb/client/v2"
	pb "github.com/mnbbrown/healthsignal/healthsignal"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

var (
	port = flag.Int("port", 10000, "HealthSignal API Port")
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
			log.Println(len(pings))
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
	grpcServer.Serve(lis)
}
