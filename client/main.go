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

type endpoint struct {
	id             int
	url            string
	expectedStatus int
}

func check(e endpoint) {
	log.Printf("Checking: %s", e.url)
	tp := newTransport()
	httpClient := &http.Client{Transport: tp}
	response, err := httpClient.Get(e.url)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Printf("Failed to get endpoint. Timed out")
		} else {
			log.Printf("Failed to get endpoint %s", err)
		}
		return
	}
	if response.StatusCode != e.expectedStatus {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ping := &pb.Ping{
		Status:             int32(response.StatusCode),
		RequestDuration:    int32(tp.ReqDuration()),
		ConnectionDuration: int32(tp.ConnDuration()),
		Location:           *location,
		TimedOut:           false,
		Endpoint:           int32(e.id),
	}
	_, err = client.SavePing(ctx, ping)
	if err != nil {
		log.Printf("Failed to save point: %v", err)
	} else {
		log.Printf("Success")
	}
}

func run() {
	endpoints := []endpoint{
		endpoint{id: 1, url: "https://dev.api.primaryledger.labrys.group/_util/_healthcheck", expectedStatus: 200},
		endpoint{id: 2, url: "https://lab.primaryledger.io/_util/_healthcheck", expectedStatus: 200},
		endpoint{id: 3, url: "https://qa.api.primaryledger.labrys.group/_util/_healthcheck", expectedStatus: 200},
	}
	for _, e := range endpoints {
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
