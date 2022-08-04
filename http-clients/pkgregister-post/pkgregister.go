package pkgregister

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// models the data we send to the HTTP server
type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// models the response when you register a new package
type pkgRegisterResult struct {
	ID string `json:"id"`
}

// url: HTTP server we send data to
// data: object that we serialize as JSON to send in request body
func registerPackageData(url string, data pkgData) (pkgRegisterResult, error) {
	// create empty struct to accept the response
	p := pkgRegisterResult{}
	// serialize the data into a JSON byte slice
	b, err := json.Marshal(data)
	if err != nil {
		return p, err
	}
	// NewReader reads from the byte slice it accepts
	// https://pkg.go.dev/bytes@go1.19#NewReader
	reader := bytes.NewReader(b)
	// send JSON byte slice to url as JSON
	r, err := http.Post(url, "application/json", reader)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()

	// send req to the server
	// ReadAll reads until error or EOF
	respData, err := io.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	if r.StatusCode != http.StatusOK {
		return p, errors.New(string(respData))
	}
	// deserialize the response and store it in
	//   the packageRegisterResult
	err = json.Unmarshal(respData, &p)
	return p, err
}
