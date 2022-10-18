package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// return a slice of bytes that represents the Response Body of an HTTP request
func fetchRemoteResource(client *http.Client, url string) ([]byte, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

// If there are any redirects from the server, return an error. Otherwise, return
// nil
// req: The request to follow the URL that the server returned
// via: Slice of requests that have been made to the server.
func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	if len(via) >= 1 {
		return errors.New(fmt.Sprintf("Attempted redirect to %s", req.URL))
	}
	return nil
}

// Returns a client that times out after d time and employs a redirect policy
func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d, CheckRedirect: redirectPolicyFunc}
	return &client
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Must specify an HTTP URL to get data from")
		os.Exit(1)
	}
	client := createHTTPClientWithTimeout(15 * time.Second)
	body, err := fetchRemoteResource(client, os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s\n", body)
}
