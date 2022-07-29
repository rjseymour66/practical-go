package main

import (
	"bytes"
	"errors"
	"testing"
)

type testConfig struct {
	args           []string
	err            error
	numTimes       int
	outputHtmlPath string
}

func TestParseArgs(t *testing.T) {
	tests := []testConfig{
		{
			args:     []string{"-h"},
			err:      errors.New("flag: help requested"),
			numTimes: 0,
		},
		{
			args:     []string{"-n", "10"},
			err:      nil,
			numTimes: 10,
		},
		{
			args:           []string{"-n", "10", "-o", "output.html"},
			err:            nil,
			numTimes:       10,
			outputHtmlPath: "output.html",
		},
		{
			args:     []string{"-n", "abc"},
			err:      errors.New("invalid value \"abc\" for flag -n: parse error"),
			numTimes: 0,
		},
		{
			args:     []string{"-n", "1", "foo"},
			err:      errors.New("Positional arguments specified"),
			numTimes: 1,
		},
	}

	// mimic os.Stdout for testing
	byteBuf := new(bytes.Buffer)

	// parseArgs(w io.Writer, args []string) (config, error)
	// 1. pass the test case method argument to the method to create the initializing values
	// 2. assign values to the return values
	// 3. Compare each possible outcome for all test case return values
	//    and what the test case args returned to the initializing value
	for _, tc := range tests {
		c, err := parseArgs(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Errorf("Expected error to be: %v, got: %v\n", tc.err, err)
		}

		if c.numTimes != tc.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v\n", tc.numTimes, c.numTimes)
		}

		if len(tc.outputHtmlPath) != 0 && c.outputHtmlPath != tc.outputHtmlPath {
			t.Errorf("Expected outputHtmlPath to be: %v, got: %v\n", tc.outputHtmlPath, c.outputHtmlPath)
		}
		byteBuf.Reset()
	}
}
