package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// curl or wget copy
func fetchRemoteResource(url string) ([]byte, error) {
	// makes GET req and returns a Response object and error value
	// https://pkg.go.dev/net/http#Response
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// Response obj has a Body method of type io.ReadCloser
	defer r.Body.Close()
	// ReadAll returns a byte slice of read values and an error
	return io.ReadAll(r.Body)
}

// The user has to pass an argument to the cmd line
// errors or the body are printed to standard out
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Must specify an HTTP URL to get data from")
		os.Exit(1)
	}
	body, err := fetchRemoteResource(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s\n", body)
}
