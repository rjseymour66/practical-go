package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestHandleHttp(t *testing.T) {
	usageMessage := `
http: A HTTP client.

http: <options> server

Options: 
  -verb string
    	HTTP method (default "GET")
`
	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		// test sub-command with no positional arg
		{
			args: []string{},
			err:  ErrNoServerSpecified,
		},
		// test sub-command when called with -h
		{
			args:   []string{"-h"},
			err:    errors.New("flag: help requested"),
			output: usageMessage,
		},
		// test sub-command with pos arg for server
		{
			args:   []string{"http://localhost"},
			err:    nil,
			output: "Executing http command\n",
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleHttp(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got: %v", err)
		}

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error %v, got %v", tc.err, err)
		}

		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Errorf("Expected output to be: %#v, got: %#v", tc.output, gotOutput)
			}
		}
		byteBuf.Reset()
	}
}
