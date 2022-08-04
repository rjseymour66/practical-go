package pkgquery

import (
	"encoding/json"
	"io"
	"net/http"
)

type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func fetchPackageData(url string) ([]pkgData, error) {
	var packages []pkgData
	// Get req to package server
	r, err := http.Get(url)
	// return empty slice and error
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// deserialize JSON data only
	// The web backend might also serve data as HTML in browser
	// https://pkg.go.dev/net/http#Header.Get
	if r.Header.Get("Content-Type") != "application/json" {
		return packages, nil
	}
	// https://pkg.go.dev/io#ReadAll
	// reads from r until error or EOF
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return packages, err
	}
	// (data to deserialize, object to deserialize into)
	err = json.Unmarshal(data, &packages)
	return packages, err
}
