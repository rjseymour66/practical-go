package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// creates HTTP test server that serves test content
func startBadTestHTTPServerV1() *httptest.Server {
	// .NewServer returns Server obj and accepts HandlerFunc
	// https://pkg.go.dev/net/http/httptest#Server
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(60 * time.Second)
				fmt.Fprintf(w, "Hello World")
			}))
	return ts
}

func TestFetchBadRemoteResourceV1(t *testing.T) {
	ts := startBadTestHTTPServerV1()
	defer ts.Close()

	client := createHTTPClientWithTimeout(200 * time.Millisecond)
	_, err := fetchRemoteResource(client, ts.URL)
	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("Expected error to contain: context deadline exceeded, Got: %v", err.Error())
	}
}
