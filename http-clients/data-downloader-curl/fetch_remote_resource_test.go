package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// creates HTTP test server that serves test content
func startTestHTTPServer() *httptest.Server {
	// .NewServer returns Server obj and accepts HandlerFunc
	// https://pkg.go.dev/net/http/httptest#Server
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello World")
			}))
	return ts
}

func TestFetchRemoteResource(t *testing.T) {
	ts := startTestHTTPServer()
	defer ts.Close()

	expected := "Hello World"

	// pass the func the Server obj URL property
	// which is a random port on localhost
	// it returns whatever the HandlerFunc tells it to do
	data, err := fetchRemoteResource(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(data) {
		t.Errorf("Expected response to be: %s, got: %s", expected, data)
	}
}
