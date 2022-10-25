package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "index",
			path:     "/api",
			expected: "Hello, world!",
		},
		{
			name:     "healthcheck",
			path:     "/healthz",
			expected: "ok",
		},
	}

	// create the server mux to test
	mux := http.NewServeMux()
	setupHandlers(mux)

	// start the test server
	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// GET req with the test server URL and the test case path
			// The test server URL is localhost + port number
			resp, err := http.Get(ts.URL + tc.path)
			// store the response body
			respBody, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			if string(respBody) != tc.expected {
				t.Errorf(
					"Expected: %s, Got: %s",
					tc.expected, string(respBody),
				)
			}
		})
	}
}
