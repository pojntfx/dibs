package utils

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// ManageableCommand is a manageable command
type ManageableCommand struct {
	execLine               string
	stdoutChan, stderrChan chan string
	dir                    string
	instance               *exec.Cmd
}

// NewManageableCommand creates a new ManageableCommand
func NewManageableCommand(execLine, dir string, stdoutChan chan string, stderrChan chan string) *ManageableCommand {
	return &ManageableCommand{
		execLine:   execLine,
		dir:        dir,
		stdoutChan: stdoutChan,
		stderrChan: stderrChan,
	}
}

func readFromReader(reader io.Reader, outChan chan string) {
	bufStdout := bufio.NewReader(reader)

	for {
		line, _, err := bufStdout.ReadLine()
		if err != nil {
			return
		}

		outChan <- string(line)
	}
}

func getCommandWrappedInSh(execLine string) *exec.Cmd {
	wrappedArgs := append([]string{"sh", "-c"}, execLine)

	return exec.Command(wrappedArgs[0], wrappedArgs[1:]...)
}

// Start starts the command
func (r *ManageableCommand) Start() error {
	r.instance = getCommandWrappedInSh(r.execLine)
	// TODO: Add test that checks if command gets executed in the set dir
	r.instance.Dir = r.dir

	stdout, err := r.instance.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := r.instance.StderrPipe()
	if err != nil {
		return err
	}

	go readFromReader(stdout, r.stdoutChan)
	go readFromReader(stderr, r.stderrChan)

	r.instance.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return r.instance.Start()
}

// Wait waits for the command to complete
func (r *ManageableCommand) Wait() error {
	if err := r.instance.Wait(); err != nil && err.Error() != "signal: killed" && err.Error() != "exec: Wait was already called" {
		return err
	}

	return nil
}

// TODO: Add test for Zombie processes
// Stop stops the command
func (r *ManageableCommand) Stop() error {
	noSuchProcessError := "no such process"

	processGroupID, err := syscall.Getpgid(r.instance.Process.Pid)
	if err != nil && err.Error() == noSuchProcessError {
		return nil
	}
	if err != nil {
		return err
	}

	// Ignore Zombie processes, which can't be killed
	// pid + 1 because we execute everything through `sh` so we have to kill it as well
	for _, pid := range []int{r.instance.Process.Pid, processGroupID, r.instance.Process.Pid + 1, processGroupID + 1} {
		err = syscall.Kill(pid, syscall.SIGKILL)
		if err != nil && err.Error() == noSuchProcessError {
			return nil
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: Add test for Zombie processes
// IsStopped returns true if the command has stopped. WARNING: This always returns false for Zombie processes.
func (r *ManageableCommand) IsStopped() bool {
	process, err := os.FindProcess(r.instance.Process.Pid)
	if err != nil {
		return true
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return true
	}

	return false
}

// GetExecLine returns the command's execLine
func (r *ManageableCommand) GetExecLine() string {
	return r.execLine
}

// Dir returns the command's dir
func (r *ManageableCommand) GetDir() string {
	return r.dir
}

// GetStdoutChan returns the command's stdout channel
func (r *ManageableCommand) GetStdoutChan() chan string {
	return r.stdoutChan
}

// GetStderrChan returns the command's stderr channel
func (r *ManageableCommand) GetStderrChan() chan string {
	return r.stderrChan
}
