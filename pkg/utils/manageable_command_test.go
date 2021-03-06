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
	testDir                            = "."
)

func TestCreateManageableCommand(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandCreate, testDir, stdoutChan, stderrChan)

	if c == nil {
		t.Error("New manageable command is nil")
	}

	if c.execLine != testCommandCreate {
		t.Error("execLine not set correctly")
	}

	if c.dir != testDir {
		t.Error("dir not set correctly")
	}

	if c.stdoutChan != stdoutChan {
		t.Error("stdoutChan not set correctly")
	}

	if c.stderrChan != stderrChan {
		t.Error("stderrChan not correctly")
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

	c := NewManageableCommand(testCommandStart, testDir, stdoutChan, stderrChan)

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

	c := NewManageableCommand(testCommandStop, testDir, stdoutChan, stderrChan)

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

	c := NewManageableCommand(testCommandIsStoppedRunningProcess, testDir, stdoutChan, stderrChan)

	defer func() {
		if err := c.Stop(); err != nil {
			t.Error(err)
		}
	}()
	if err := c.Start(); err != nil {
		t.Error(err)
	}

	if processIsStopped := c.IsStopped(); processIsStopped != false {
		t.Error("command is running but IsStopped returned false")
	}
}

func TestIsStoppedStoppedProcess(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandIsStoppedStoppedProcess, testDir, stdoutChan, stderrChan)

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

	c := NewManageableCommand(testCommandGetters, testDir, stdoutChan, stderrChan)

	if c.GetExecLine() != testCommandGetters {
		t.Error("GetExecLine did not return the set execLine")
	}
}

func TestGetDir(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandGetters, testDir, stdoutChan, stderrChan)

	if c.GetDir() != testDir {
		t.Error("GetDir did not return the set dir")
	}
}

func TestGetStdoutChan(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandGetters, testDir, stdoutChan, stderrChan)

	if c.GetStdoutChan() == nil {
		t.Error("GetStdoutChan has not been set")
	}
}

func TestGetStderrChan(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	c := NewManageableCommand(testCommandGetters, testDir, stdoutChan, stderrChan)

	if c.GetStderrChan() == nil {
		t.Error("GetStderrChan has not been set")
	}
}
