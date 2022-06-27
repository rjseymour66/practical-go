package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runCmd(r io.Reader, w io.Writer, c config) error
func TestRunCmd(t *testing.T) {
	tests := []struct {
		c            config
		input        string
		output       string
		fileContents string
		err          error
	}{
		{
			c:      config{numTimes: 5},
			input:  "",
			output: strings.Repeat("Your name please? Press the Enter key when done.\n", 1),
			err:    errors.New("You didn't enter your name"),
		},
		{
			c:      config{numTimes: 5},
			input:  "Bill Bryson",
			output: "Your name please? Press the Enter key when done.\n" + strings.Repeat("Nice to meet you Bill Bryson\n", 5),
		},
		{
			c:      config{numTimes: 5, name: "Bill Bryson"},
			input:  "",
			output: strings.Repeat("Nice to meet you Bill Bryson\n", 5),
		},
		{
			c:            config{numTimes: 0, outputHtmlPath: "output.html"},
			input:        "Jane Clancy",
			output:       "Your name please? Press the Enter key when done.\n",
			fileContents: "<h1>Hello Jane Clancy</h1>",
		},
	}

	// test the io.Writer (output)
	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		// creates an io.Reader from a string
		rd := strings.NewReader(tc.input)

		if len(tc.c.outputHtmlPath) != 0 {
			tc.c.outputHtmlPath = filepath.Join(t.TempDir(), tc.c.outputHtmlPath)
		}

		err := runCmd(rd, byteBuf, tc.c)
		if err != nil && tc.err == nil {
			t.Fatalf("Expected nil error, got: %v\n", err)
		}

		if tc.err != nil {
			if err.Error() != tc.err.Error() {
				t.Fatalf("Expected error: %v, Got error: %v\n", tc.err.Error(), err.Error())
			}
		}

		// convert contents of byteBuf to String()
		gotMsg := byteBuf.String()
		if gotMsg != tc.output {
			t.Errorf("Expected stdout message to be: %v, Got: %v\n", tc.output, gotMsg)
		}

		if len(tc.fileContents) != 0 {
			fileData, err := os.ReadFile(tc.c.outputHtmlPath)
			if err != nil {
				t.Fatal(err)
			}
			if string(fileData) != tc.fileContents {
				t.Errorf("Expected file contents to be: %v, got: %v\n", tc.fileContents, string(fileData))
			}
		}

		// empties the buffer before loading the next test case
		byteBuf.Reset()
	}
}
