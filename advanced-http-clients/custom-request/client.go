package main

import (
	"context"
	"net/http"
)

// When you make a request with a Client object, you use the default Request object.
// You can create a custom request object and add context to it.
func createHTTPGetRequest(ctx context.Context, url string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req, err
}
