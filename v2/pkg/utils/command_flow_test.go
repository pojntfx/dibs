package utils

import (
	"strings"
	"testing"
)

var (
	testCommandsCreate = []string{"ls", "ls -la"}
	testCommandsStart  = []string{"ping -c 1 localhost", "ping -c 1 127.0.0.1"}
)

func TestCreateCommandFlow(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	f := NewCommandFlow(testCommandsCreate, stdoutChan, stderrChan)

	if len(f.commands) < 2 {
		t.Error("commands not set")
	}
}

func TestStartCommandFlow(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	hits := 0
	go func() {
		for {
			select {
			case stdout := <-stdoutChan:
				t.Log("test stdout", stdout)

				if strings.Contains(stdout, "PING localhost") {
					hits++
				}

				if strings.Contains(stdout, "PING 127.0.0.1") {
					hits++
				}
			case stderr := <-stderrChan:
				t.Error("error while executing command", stderr)
			}
		}
	}()

	f := NewCommandFlow(testCommandsStart, stdoutChan, stderrChan)

	if err := f.Start(); err != nil {
		t.Error(err)
	}

	if hits < 2 {
		t.Error("commands did not match expected output")
	}
}
