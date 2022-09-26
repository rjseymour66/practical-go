package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Accepts an HTTP client and URL to fetch data from.
func fetchRemoteResource(client *http.Client, url string) ([]byte, error) {
	// Use the client's Get method instead of the http.Get method. http.client
	// has useful methods and attributes to create custom clients using Go's
	// standard library.
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

// creates a client with the specified timeout duration. We could create a client
// that has any customized methods or attributes and pass it to
// fetchRemoteResource to return data.
func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Must specify an HTTP URL to get data from")
		os.Exit(1)
	}
	client := createHTTPClientWithTimeout(15 * time.Second)
	body, err := fetchRemoteResource(client, os.Args[1])
	// Fprintf formats a string message and writes to a w Writer (such as STDOUT)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%#v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s\n", body)
}
