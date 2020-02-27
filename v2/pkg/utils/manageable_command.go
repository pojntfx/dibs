package utils

import (
	"bufio"
	"io"
	"os/exec"
	"syscall"
)

// ManageableCommand is a manageable command
type ManageableCommand struct {
	execLine               string
	stdoutChan, stderrChan chan string
	instance               *exec.Cmd
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
func (r *ManageableCommand) Wait() error { return r.instance.Wait() }

// Stop stops the command
func (r *ManageableCommand) Stop() error {
	processGroupID, err := syscall.Getpgid(r.instance.Process.Pid)
	if err != nil {
		return err
	}

	return syscall.Kill(processGroupID, syscall.SIGKILL)
}

// NewManageableCommand creates a new ManageableCommand
func NewManageableCommand(execLine string, stdoutChan chan string, stderrChan chan string) *ManageableCommand {
	return &ManageableCommand{
		execLine:   execLine,
		stdoutChan: stdoutChan,
		stderrChan: stderrChan,
	}
}
