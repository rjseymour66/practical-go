package pkgregister

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type pkgData struct {
	Name     string
	Version  string
	Filename string
	Bytes    io.Reader // reads files
}

type pkgRegisterResult struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func registerPackageData(url string, data pkgData) (pkgRegisterResult, error) {

	// create empty struct to store response
	p := pkgRegisterResult{}
	// create the binary payload
	payload, contentType, err := createMultiPartMessage(data)
	if err != nil {
		return p, err
	}

	// store the payload in reader
	reader := bytes.NewReader(payload)
	// send the POST, store in r
	r, err := http.Post(url, contentType, reader)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()
	// read the response
	respData, err := io.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	// store the response in the empty struct
	err = json.Unmarshal(respData, &p)
	return p, err
}

func createHTTPCLientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}
