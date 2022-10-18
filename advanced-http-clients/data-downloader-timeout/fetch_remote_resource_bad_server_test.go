package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Using a channel to test a remote resource.
func startBadTestHTTPServerV2(shutdownServer chan struct{}) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-shutdownServer
		fmt.Fprint(w, "Hello World")
	}))
	return ts
}

func TestFetchBadRemoteResourceV2(t *testing.T) {
	// unbuffered channel of type empty struct{}
	shutdownServer := make(chan struct{})
	ts := startBadTestHTTPServerV2(shutdownServer)
	defer ts.Close()
	// writes an empty struct to the channel before ts.Close() is called
	defer func() {
		shutdownServer <- struct{}{}
	}()

	client := createHTTPClientWithTimeout(200 * time.Millisecond)
	_, err := fetchRemoteResource(client, ts.URL)
	if err == nil {
		t.Log("Expected non-nil error")
		t.Fail()
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("Expected error to contain: context deadline exceeded, Got: %v", err.Error())
	}
}
