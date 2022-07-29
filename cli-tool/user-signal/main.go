package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func createContextWithTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	return ctx, cancel
}

// creates a channel that calls a cancelling function when one of the signals is received
func setupSignalHandler(w io.Writer, cancelFunc context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-c
		fmt.Fprintf(w, "Got signal: %v\n", s)
		cancelFunc()
	}()
}

func executeCommand(ctx context.Context, command string, arg string) error {
	// CommandContext force kills an external program with the context expires
	return exec.CommandContext(ctx, command, arg).Run()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stdout, "Usage: %s <command> <argument>\n", os.Args[0])
		os.Exit(1)
	}
	command := os.Args[1]
	arg := os.Args[2]

	// Step 1: create a context using WithTimeout()
	cmdTimeout := 30 * time.Second
	ctx, cancel := createContextWithTimeout(cmdTimeout)
	defer cancel()

	// Step 2: Call cancelfunction when one of the signals is recieved
	setupSignalHandler(os.Stdout, cancel)

	// Step 3: Execute the external program using the CommandContext() func
	// 		   using context from Step 1
	err := executeCommand(ctx, command, arg)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
