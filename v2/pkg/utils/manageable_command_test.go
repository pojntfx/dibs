package utils

import (
	"strings"
	"testing"
)

func TestStart(t *testing.T) {
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

	c := New("ping -c 1 localhost", stdoutChan, stderrChan)

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
