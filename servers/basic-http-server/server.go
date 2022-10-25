package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// check for env var
	listenAddr := os.Getenv("LISTEN_ADDR")
	// if env var is not set, listen on 8080 on any machine
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	// takes addr and handler.
	// returns with error if started incorrectly or when terminated
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
