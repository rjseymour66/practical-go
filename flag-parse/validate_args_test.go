package main

import (
	"errors"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		c   config
		err error
	}{
		{
			c:   config{},
			err: errors.New("Must specify a number greater than 0"),
		},
		{
			c:   config{outputHtmlPath: "output.html"},
			err: nil,
		},
		{
			c:   config{numTimes: -1},
			err: errors.New("Must specify a number greater than 0"),
		},
		{
			c:   config{numTimes: 10},
			err: nil,
		},
	}

	// loop through the slice of test cases (tc)
	// assign the return type to a variable
	// if clauses to see if return type
	for _, tc := range tests {
		err := validateArgs(tc.c)
		// err.Error() == fmt.Print(err)
		// if the test case returns an error, and the validateArgs
		//    error doesn't match the tc error
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Errorf("Expected error to be: %v, got: %v\n", tc.err, err)
		}

		// if the test case does not return an error and validateArgs does...
		if tc.err == nil && err != nil {
			t.Errorf("Expected nil error, got: %v\n", err)
		}
	}
}
