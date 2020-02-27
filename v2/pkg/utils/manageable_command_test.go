package utils

import (
	"strings"
	"testing"
)

const (
	testCommandCreate = "ls"
	testCommandStart  = "ping -c 1 localhost"
)

func TestCreateManageableCommand(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandCreate, stdoutChan, stderrChan)

	if c.execLine == "" {
		t.Error("exec line not set")
	}

	if c.stdoutChan == nil {
		t.Error("stdoutChan not set")
	}

	if c.stderrChan == nil {
		t.Error("stderrChan not set")
	}
}

func TestStartManageableCommand(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	hits := 0
	go func() {
		for {
			select {
			case stdout := <-stdoutChan:
				t.Log("test stdout", stdout)

				if strings.Contains(stdout, "PING") {
					hits++
				}
			case stderr := <-stderrChan:
				t.Error("error while executing command", stderr)
			}
		}
	}()

	c := NewManageableCommand(testCommandStart, stdoutChan, stderrChan)

	if err := c.Start(); err != nil {
		t.Error(err)
	}

	if err := c.Wait(); err != nil {
		t.Error(err)
	}

	if hits < 0 {
		t.Error("command did not match expected output")
	}
}
