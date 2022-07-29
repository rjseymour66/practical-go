package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// timeout after 5 seconds
var totalDuration time.Duration = 5

// reads a name from stdin
func getName(r io.Reader, w io.Writer) (string, error) {
	scanner := bufio.NewScanner(r)
	msg := "Your name please? Press the Enter key when done"
	fmt.Fprintln(w, msg)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("You entered an empty name")
	}
	return name, nil
}

// returns the default name and an error, or the
// getName return values
func getNameContext(ctx context.Context, r io.Reader, w io.Writer) (string, error) {
	var err error
	name := "Default name"
	// buffered channels define a length
	c := make(chan error, 1)

	//
	go func() {
		name, err = getName(r, w)
		// write the error in the buffered channel
		c <- err
	}()

	// blocks until one of the cases can run, then runs
	// the case
	select {
	// if the ctx expires, return the default name
	// and the error
	case <-ctx.Done():
		return name, ctx.Err()
		// receive the error from c and assign to err
	case err := <-c:
		return name, err
	}
}

func main() {
	allowedDuration := totalDuration * time.Second

	// create parent context,
	ctx, cancel := context.WithTimeout(context.Background(), allowedDuration)
	defer cancel()

	name, err := getNameContext(ctx, os.Stdin, os.Stdout)
	// if gnc returns bc of expired ctx with error, return the default name
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}
	// else, return the entered name to Stdout
	fmt.Fprintln(os.Stdout, name)
}
