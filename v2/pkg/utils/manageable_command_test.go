package utils

import (
	"strings"
	"testing"
	"time"
)

const (
	testCommandCreate                  = "ls"
	testCommandStart                   = "ping -c 1 localhost"
	testCommandStop                    = "ping -c 60 localhost"
	testCommandIsStoppedRunningProcess = testCommandStop
	testCommandIsStoppedStoppedProcess = testCommandStop
	testCommandGetters                 = testCommandCreate
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

				if strings.Contains(stdout, "PING localhost") {
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

	if hits <= 0 {
		t.Error("command did not match expected output")
	}
}

func TestStopManageableCommand(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandStop, stdoutChan, stderrChan)

	if err := c.Start(); err != nil {
		t.Error(err)
	}

	go func() {
		time.Sleep(time.Second * 1)

		if err := c.Stop(); err != nil {
			t.Error(err)
		}
	}()

	if err := c.Wait(); err != nil {
		t.Error(err)
	}
}

func TestIsStoppedRunningProcess(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandIsStoppedRunningProcess, stdoutChan, stderrChan)

	defer c.Stop()
	if err := c.Start(); err != nil {
		t.Error(err)
	}

	if processIsStopped := c.IsStopped(); processIsStopped != false {
		t.Error("command is running but IsStopped returned false")
	}
}

func TestIsStoppedStoppedProcess(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandIsStoppedStoppedProcess, stdoutChan, stderrChan)

	if err := c.Start(); err != nil {
		t.Error(err)
	}

	if err := c.Stop(); err != nil {
		t.Error(err)
	}

	if err := c.Wait(); err != nil {
		t.Error(err)
	}

	if processIsStopped := c.IsStopped(); processIsStopped != true {
		t.Error("command is not running but IsStopped returned true")
	}
}

func TestGetExecLine(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandGetters, stdoutChan, stderrChan)

	if c.GetExecLine() != testCommandGetters {
		t.Error("GetExecLine did not return the set execLine")
	}
}

func TestGetStdoutChan(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandGetters, stdoutChan, stderrChan)

	if c.GetStdoutChan() == nil {
		t.Error("GetStdoutChan has not been set")
	}
}

func TestGetStderrChan(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandGetters, stdoutChan, stderrChan)

	if c.GetStderrChan() == nil {
		t.Error("GetStderrChan has not been set")
	}
}
