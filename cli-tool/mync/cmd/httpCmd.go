package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type httpConfig struct {
	url  string
	verb string
}

func validateConfig(c httpConfig) error {
	allowedVerbs := []string{"GET", "POST", "HEAD"}
	for _, v := range allowedVerbs {
		if c.verb == v {
			return nil
		}
	}
	return ErrInvalidHTTPMethod
}

func HandleHttp(w io.Writer, args []string) error {
	c := httpConfig{}

	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	// sets the default
	fs.StringVar(&c.verb, "verb", "GET", "HTTP method")

	fs.Usage = func() {
		var usageString = `
http: A HTTP client.

http: <options> server`
		fmt.Fprintf(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	err = validateConfig(c)
	if err != nil {
		if errors.Is(err, ErrInvalidHTTPMethod) {
			fmt.Fprintln(w, "Invalid HTTP method")
		}
		return err
	}

	c.url = fs.Arg(0)
	data, err := fetchRemoteResource(c.url)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, string(data))
	return nil
}

// curl or wget copy
func fetchRemoteResource(url string) ([]byte, error) {
	// makes GET req and returns a Response object and error value
	// https://pkg.go.dev/net/http#Response
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// Response obj has a Body method of type io.ReadCloser
	defer r.Body.Close()
	// ReadAll returns a byte slice of read values and an error
	return io.ReadAll(r.Body)
}
