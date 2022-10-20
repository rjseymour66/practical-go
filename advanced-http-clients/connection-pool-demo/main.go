package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}

// httptrace lets you see if a new connection is being established or if you are
// reusing an existing one.
// Accepts a context and URL string.
// Returns a request with a new context that traces TCP connections.
func createHTTPGetRequestWithTrace(ctx context.Context, url string) (*http.Request, error) {
	// create a req
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	// ClientTrace defines functions that are called when a specific event in a
	// connections life cycle occurs.
	trace := &httptrace.ClientTrace{
		// When DNS lookup for a hostname event is complete, it calls the associated
		// function. The func must use this param and cannot return a value.
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
		// When the connection is obtained event is complete, it calls the associated
		// function. The func must use this param and cannot return a value.
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
	}
	// Create new context that traces connection events
	ctxTrace := httptrace.WithClientTrace(req.Context(), trace)
	// Add the newly created context
	req = req.WithContext(ctxTrace)
	return req, err
}

func main() {
	d := 5 * time.Second
	ctx := context.Background()
	client := createHTTPClientWithTimeout(d)

	req, err := createHTTPGetRequestWithTrace(ctx, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for {
		client.Do(req)
		time.Sleep(1 * time.Second)
		fmt.Println("---------")
	}
}
