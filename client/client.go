package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	pb "github.com/mnbbrown/healthsignal/healthsignal"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

var client pb.HealthSignalClient

var (
	serverAddr = flag.String("server_addr", "grpc.healthsignal.live:10443", "The server address in the format of host:port")
	location   = flag.String("location", "london", "The location we're checking from")
	useTLS     = flag.Bool("tls", true, "Connect with TLS")
)

func createBody(body string) io.Reader {
	return strings.NewReader(body)
}

func newRequest(method string, url *url.URL, body string) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), createBody(body))
	return req, err
}

func dialContext(network string) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, _, addr string) (net.Conn, error) {
		return (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
		}).DialContext(ctx, network, addr)
	}
}

func readResponseBody(req *http.Request, resp *http.Response) error {
	w := ioutil.Discard
	if _, err := io.Copy(w, resp.Body); err != nil && w != ioutil.Discard {
		return err
	}
	return nil
}

func check(e *pb.Endpoint) {
	log.Printf("Checking: %s", e.Url)

	url, err := url.Parse(e.Url)
	if err != nil {
		log.Printf("Failed to parse URL: %v", err)
		return
	}

	req, err := newRequest(e.Method, url, "")
	if err != nil {
		log.Printf("Failed to build request: %v", err)
		return
	}

	var dnsStart, dnsDone, connStart, connDone, gotConn, gotFirstByte, tlsStart, tlsEnd, bodyReadEnd time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { dnsDone = time.Now() },
		ConnectStart: func(_, _ string) {
			connStart = time.Now()
		},
		ConnectDone: func(net, addr string, err error) {
			if err != nil {

			}
			connDone = time.Now()
		},
		GotConn:              func(_ httptrace.GotConnInfo) { gotConn = time.Now() },
		GotFirstResponseByte: func() { gotFirstByte = time.Now() },
		TLSHandshakeStart:    func() { tlsStart = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { tlsEnd = time.Now() },
	}
	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	tr.DialContext = dialContext("tcp4")

	switch url.Scheme {
	case "https":
		host, _, err := net.SplitHostPort(req.Host)
		if err != nil {
			host = req.Host
		}

		tr.TLSClientConfig = &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: true,
		}

		err = http2.ConfigureTransport(tr)
		if err != nil {
			log.Fatalf("failed to prepare transport for HTTP/2: %v", err)
		}
	}

	httpClient := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("failed to read response: %v", err)
		return
	}
	err = readResponseBody(req, resp)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("failed to read body: %v", err)
		return
	}
	bodyReadEnd = time.Now()

	onlineStatus := true
	if int32(resp.StatusCode) != e.ExpectedStatus {
		onlineStatus = false
	}

	deltaToMs := func(start, end time.Time) int32 {
		return int32(end.Sub(start) / time.Millisecond)
	}

	// Send to server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ping := &pb.Ping{
		Endpoint:                 int32(e.Id),
		Location:                 *location,
		HttpStatus:               int32(resp.StatusCode),
		Protocol:                 "http",
		DnsLookupDuration:        deltaToMs(dnsStart, dnsDone),
		TcpConnectionDuration:    deltaToMs(connStart, connDone),
		TlsHandshakeDuration:     deltaToMs(tlsStart, tlsEnd),
		ServerProcessingDuration: deltaToMs(gotConn, gotFirstByte),
		ContentTransferDuration:  deltaToMs(gotFirstByte, bodyReadEnd),
		Online:                   onlineStatus,
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
		log.Printf("Failed to get endpoints %v", err)
		return
	}
	log.Printf("Checking %d endpoints", len(endpoints.Endpoints))
	for _, e := range endpoints.Endpoints {
		go check(e)
	}
}

func main() {
	flag.Parse()

	var conn *grpc.ClientConn
	if *useTLS {
		pool, err := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		conn, err = grpc.Dial(*serverAddr, grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Fatalf("Failed to connect via tls: %v", err)
		}
		log.Printf("Successfully connected to %s with TLS", *serverAddr)
	} else {
		var err error
		conn, err = grpc.Dial(*serverAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		log.Printf("Successfully connected to %s without TLS", *serverAddr)
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
