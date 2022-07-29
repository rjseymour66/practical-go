package main

import (
	"errors"
	"testing"
)

type testConfig struct {
	args []string
	err  error
	config
}

func TestParseArgs(t *testing.T) {
	tests := []testConfig{
		{
			args:   []string{"-h"},
			err:    nil,
			config: config{printUsage: true, numTimes: 0},
		},
		{
			args:   []string{"--help"},
			err:    nil,
			config: config{printUsage: true, numTimes: 0},
		},
		{
			args:   []string{"10"},
			err:    nil,
			config: config{printUsage: false, numTimes: 10},
		},
		{
			args:   []string{"abc"},
			err:    errors.New("strconv.Atoi: parsing \"abc\": invalid syntax"),
			config: config{printUsage: false, numTimes: 0},
		},
		{
			args:   []string{"1", "foo"},
			err:    errors.New("Invalid number of arguments"),
			config: config{printUsage: false, numTimes: 0},
		},
	}

	// parseArgs(args []string) (config, error)
	// 1. pass the test case method argument to the method to create the initializing values
	// 2. assign values to the return values
	// 3. Compare each possible outcome for all test case return values
	//       and what the test case args returned to the initializing value
	for _, tc := range tests {
		c, err := parseArgs(tc.args)

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error to be: %v, got: %v\n", tc.err, err)
		}

		if tc.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}

		if c.printUsage != tc.printUsage {
			t.Errorf("Expected printUsage to be: %v, got: %v\n", tc.printUsage, c.printUsage)
		}

		if c.numTimes != tc.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v\n", tc.numTimes, c.numTimes)
		}
	}
}
