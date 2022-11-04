package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

type requestContextKey struct{}
type requestContextValue struct {
	requestID string
}

// helper that stores the req identifer in the req's context
func addRequestID(r *http.Request, requestID string) *http.Request {
	c := requestContextValue{
		requestID: requestID,
	}
	// get current ctx
	currentCtx := r.Context()
	// create new context by adding the requestID to the existing context
	newCtx := context.WithValue(currentCtx, requestContextKey{}, c)
	// return existing context with the new parts of the new ctx
	return r.WithContext(newCtx)
}

// logs a valid requestID from the context
func logRequest(r *http.Request) {
	ctx := r.Context()
	v := ctx.Value(requestContextKey{})

	// type assertion?
	if m, ok := v.(requestContextValue); ok {
		log.Printf("Processing request: %s", m.requestID)
	}
}

// logs the request and write that it was processed
func processRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	fmt.Fprintf(w, "Request processed")
}

// handler that creates a requestID, creates a request with a context
// that contains the requestID, then logs the request
func apiHandler(w http.ResponseWriter, r *http.Request) {
	requestID := "request-123-abc"
	r = addRequestID(r, requestID)
	processRequest(w, r)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api", apiHandler)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
