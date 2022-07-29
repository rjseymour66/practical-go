package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

var binaryName string

// build then delete binary
func TestMain(m *testing.M) {
	if runtime.GOOS == "windows" {
		binaryName = "sub-cmd-arch.exe"
	} else {
		binaryName = "sub-cmd-arch"
	}

	// only run for 5 seconds then cancel
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// build the binary
	cmd := exec.CommandContext(ctx, "go", "build", "-o", binaryName)
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}

	// cleanup binary
	defer func() {
		err = os.Remove(binaryName)
		if err != nil {
			log.Fatalf("Error removing build binary: %v", err)
		}
	}()
	m.Run()
}

// the next func tests the possible CLI args (-i), possible input ("My name"), exit code, and expected output

// func TestApplication(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)

// 	// get cwd to create binaryPath
// 	curDir, err := os.Getwd()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	binaryPath := path.Join(curDir, binaryName)
// 	t.Log(binaryPath)

// 	tests := []struct {
// 		args
// 	}
// }
