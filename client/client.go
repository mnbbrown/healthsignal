package main

import (
	"flag"
	pb "github.com/mnbbrown/healthsignal/healthsignal"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

var client pb.HealthSignalClient

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	location   = flag.String("location", "london", "The location we're checking from")
)

func check(e *pb.Endpoint) {
	log.Printf("Checking: %s", e.Url)
	tp := newTransport()
	httpClient := &http.Client{Transport: tp}
	response, err := httpClient.Get(e.Url)

	timedOut := false
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Printf("Failed to get endpoint. Timed out")
			timedOut = true
		} else {
			log.Printf("Failed to get endpoint %s", err)
			return
		}
	}
	var onlineStatus int32
	if int32(response.StatusCode) != e.ExpectedStatus {
		onlineStatus = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ping := &pb.Ping{
		Status:             int32(response.StatusCode),
		RequestDuration:    int32(tp.ReqDuration() / time.Millisecond),
		ConnectionDuration: int32(tp.ConnDuration() / time.Millisecond),
		Location:           *location,
		TimedOut:           timedOut,
		Endpoint:           int32(e.Id),
		OnlineStatus:       onlineStatus,
	}
	_, err = client.SavePing(ctx, ping)
	if err != nil {
		log.Printf("Failed to save point: %v", err)
	} else {
		log.Printf("Success")
	}
}

func run() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	endpoints, err := client.GetEndpoints(ctx, &pb.EndpointsQuery{})
	if err != nil {
		return
	}
	for _, e := range endpoints.Endpoints {
		go check(e)
	}
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewHealthSignalClient(conn)

	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			run()
		}
	}
}
