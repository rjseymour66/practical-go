package main

import (
	"context"
	"net/http"
	"time"
)

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}

func main() {
	client := createHTTPClientWithTimeout(20 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()

	req, err := createHTTPGetRequest(ctx, ts.URL+"/api/packages", nil)
	res, err := client.Do(req)
}
