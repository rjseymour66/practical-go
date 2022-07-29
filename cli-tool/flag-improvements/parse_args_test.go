package main

import (
	"bytes"
	"errors"
	"testing"
)

type testConfig struct {
	args []string
	config
	output         string
	err            error
	outputHtmlPath string
}

func TestParseArgs(t *testing.T) {
	tests := []testConfig{
		// {
		// 	args: []string{"-h"},
		// 	output: `
		// 	A greeter application that prints the name you entered a specified number of times.

		// 	Usage of greeter: <options> [name]

		// 	Options:
		// 	   -n int
		// 	   			Number of times to greet
		// 	`,
		// 	err:    errors.New("flag: help requested"),
		// 	config: config{numTimes: 0},
		// },
		{
			args:   []string{"-n", "10"},
			err:    nil,
			config: config{numTimes: 10},
		},
		{
			args:           []string{"-n", "10", "-o", "output.html"},
			err:            nil,
			config:         config{numTimes: 10},
			outputHtmlPath: "output.html",
		},
		{
			args:   []string{"-n", "abc"},
			err:    errors.New("invalid value \"abc\" for flag -n: parse error"),
			config: config{numTimes: 0},
		},
		{
			args:   []string{"-n", "1", "John Doe"},
			err:    nil,
			config: config{numTimes: 1, name: "John Doe"},
		},
		{
			args:   []string{"-n", "1", "John", "Doe"},
			err:    errors.New("More than one positional argument specified"),
			config: config{numTimes: 1},
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
			t.Fatalf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error to be: %v, got: %v\n", tc.err, err)
		}

		if c.numTimes != tc.numTimes {
			t.Fatalf("Expected numTimes to be: %v, got: %v\n", tc.numTimes, c.numTimes)
		}
		gotMsg := byteBuf.String()

		if len(tc.output) != 0 && gotMsg != tc.output {
			t.Errorf("Expected stdout message to be: %#v, Got: %#v\n", tc.output, gotMsg)
		}

		if len(tc.outputHtmlPath) != 0 && c.outputHtmlPath != tc.outputHtmlPath {
			t.Fatalf("Expected outputHtmlPath to be: %v, got: %v\n", tc.outputHtmlPath, c.outputHtmlPath)
		}
		byteBuf.Reset()
	}
}
