package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type LoggingClient struct {
	log *log.Logger
}

// Use our own typ that implements the RoundTripper interface:
// Roundtrip(*Request) (*Response, error)

// Use the logging client type to implement the RoundTrip interface. If
// you implement the RoundTrip interface, you can create a custom Transport field
// for a client
func (c LoggingClient) RoundTrip(r *http.Request) (*http.Response, error) {
	c.log.Printf("Sending a %s request to %s over %s\n", r.Method, r.URL, r.Proto)
	// use DefaultTransport to send an HTTP req. It is of type RoundTripper,
	// so it has a Roundtrip method.
	resp, err := http.DefaultTransport.RoundTrip(r)
	c.log.Printf("Got back a response over %s\n", resp.Proto)

	return resp, err
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}

func fetchRemoteResource(client *http.Client, url string) ([]byte, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

// The main function does the following:
//   1. Creates a LoggingClient struct. Because it implements the Roundtripper
//      interface, you can set as the Transport method on an HTTP client.
//   2. Creates a new logger.
//   3. Assigns the new logger to the log property of the LoggingClient struct
//   4. Creates a client with a 15 sec timeout
//   5. Assigns the LoggingClient the HTTP client's Transport method.
//   6.

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "must specify an HTTP URL to get data from")
		os.Exit(1)
	}

	// create a new LoggingClient
	myTransport := LoggingClient{}
	// create new logger that logs to stdout, prefixes each string with
	//   nothing (""), then prefix each log line with date and time
	l := log.New(os.Stdout, "", log.LstdFlags)
	// assign the new logger to the l field in the LoggingClient
	myTransport.log = l

	// create the client
	client := createHTTPClientWithTimeout(15 * time.Second)
	// LogginClient implements the RoundTripper interface, so you can assign it
	//   to the Transport method. The Transport method is what sends a req to a
	//   server, then returns a response to the client:
	//   RoundTrip(*Request) (*Response, error)
	client.Transport = &myTransport

	body, err := fetchRemoteResource(client, os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%#v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Bytes in response: %d\n", len(body))
}
