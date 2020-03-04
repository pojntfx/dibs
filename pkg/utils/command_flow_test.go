package utils

import (
	"log"
	"strings"
	"testing"
	"time"
)

var (
	testCommandsCreate  = []string{"ls", "ls -la"}
	testCommandsStart   = []string{"ping -c 1 localhost", "ping -c 1 127.0.0.1"}
	testCommandsStop    = []string{"ping -c 1 localhost", "ping -c 60 127.0.0.1"}
	testCommandsRestart = testCommandsStop
)

func TestCreateCommandFlow(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	f := NewCommandFlow(testCommandsCreate, testDir, stdoutChan, stderrChan)

	if f == nil {
		t.Error("New command flow is nil")
	}

	if len(f.commands) != len(testCommandCreate) {
		t.Error("commands not set correctly")
	}

	if f.isRestart != false {
		t.Error("isRestart not false")
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

	f := NewCommandFlow(testCommandsStart, testDir, stdoutChan, stderrChan)

	if err := f.Start(); err != nil {
		t.Error(err)
	}

	if err := f.Wait(); err != nil {
		t.Error(err)
	}

	if hits < 2 {
		t.Error("commands did not match expected output")
	}
}

func TestStopCommandFlow(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	f := NewCommandFlow(testCommandsStop, testDir, stdoutChan, stderrChan)

	if err := f.Start(); err != nil {
		t.Error(err)
	}

	go func(t *testing.T) {
		time.Sleep(time.Second * 2)

		if err := f.Stop(); err != nil {
			t.Error(err)

			log.Fatal(err)
		}
	}(t)

	if err := f.Wait(); err != nil {
		t.Error(err)
	}
}

func TestRestartCommandFlow(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	f := NewCommandFlow(testCommandsRestart, testDir, stdoutChan, stderrChan)

	if err := f.Start(); err != nil {
		t.Error(err)
	}

	go func(t *testing.T) {
		time.Sleep(time.Second * 2)

		if err := f.Restart(); err != nil {
			t.Error(err)

			log.Fatal(err)
		}

		if err := f.Stop(); err != nil {
			t.Error(err)

			log.Fatal(err)
		}
	}(t)

	if err := f.Wait(); err != nil {
		t.Error(err)
	}
}
