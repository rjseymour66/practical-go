package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type logLine struct {
	UserIP string `json:"user_ip"`
	Event  string `json:"event"`
}

// decodes individual objects when they are sent as part of
// 	the same request.
func decodeHandler(w http.ResponseWriter, r *http.Request) {
	// NewDecoder implements io.Reader, so can read the entire body
	dec := json.NewDecoder(r.Body)

	var e *json.UnmarshalTypeError

	// loop and log object to console. Stop at unknown error, or EOF and
	// send OK response to client
	for {
		// store single decoded entry in this object
		var l logLine
		// read until first valid JSON object
		// return err on invalid JSON syntax or data type
		err := dec.Decode(&l)
		// stop on EOF
		if err == io.EOF {
			break
		}
		// continue processing on unmarshalling error
		if errors.As(err, &e) {
			log.Println(err)
			continue
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(l.UserIP, l.Event)
	}
	// send OK response to the client
	fmt.Fprintf(w, "OK")
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/decode", decodeHandler)

	http.ListenAndServe(":8080", mux)
}
