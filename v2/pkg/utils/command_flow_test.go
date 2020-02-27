package utils

import "testing"

var (
	testCommandsCreate = []string{"ls", "ls -la"}
)

func TestCreate(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	f := NewCommandFlow(testCommandsCreate, stdoutChan, stderrChan)

	if len(f.commands) < 2 {
		t.Error("commands not set")
	}
}
