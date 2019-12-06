package utils

import (
	"github.com/otiai10/copy"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"os"
	"os/exec"
	"strings"
)

// EventedCommand is a command that also publishes an event
type EventedCommand struct {
	LogMessage   string // The message to log to stdout
	ExecLine     string // The command to run/start
	RedisSuffix  string // The Redis channel suffix to use for the event
	RedisMessage string // The Redis message to send
}

// Pipeline is a development configuration
type Pipeline struct {
	Module                  string           // The module that is being pushed
	ModulePushedRedisSuffix string           // The Redis suffix channel suffix to use for the pushed event
	SrcDir                  string           // The directory of the module's source
	PushDir                 string           // The temporary directory to use for the modification of the Git repo
	RunCommands             []EventedCommand // The commands to run
	StartCommand            EventedCommand   // The command to start (will keep running, but can be killed)
	StartCommandState       *exec.Cmd        // Stores the state of the start command (to make it possible to kill it)
	Git                     Git              // Git instance to use to modify the Git repo in PushDir
	Redis                   Redis            // Redis instance to use to publish the events
}

// RunCommandsOnly only runs the commands
func (pipeline *Pipeline) RunCommandsOnly() error {
	for _, runCommand := range pipeline.RunCommands {
		log.Info(runCommand.LogMessage, rz.String("Module", pipeline.Module))
		if err := RunCommand(runCommand.ExecLine, false); err != nil {
			return err
		}
		pipeline.Redis.PublishWithTimestamp(runCommand.RedisSuffix, runCommand.RedisMessage)
	}

	return nil
}

// RunAll runs the entire pipeline
func (pipeline *Pipeline) RunAll() error {
	if pipeline.StartCommandState != nil {
		log.Info("Stopping module ...", rz.String("Module", pipeline.Module))
		pipeline.StartCommandState.Process.Kill()
	}

	if err := setupPushDir(pipeline.SrcDir, pipeline.PushDir); err != nil {
		return err
	}

	if err := pipeline.Git.PushModule(pipeline.Module, pipeline.PushDir); err != nil {
		return err
	}
	pipeline.Redis.PublishWithTimestamp(pipeline.ModulePushedRedisSuffix, pipeline.Module)

	pipeline.RunCommandsOnly()

	log.Info(pipeline.StartCommand.LogMessage, rz.String("Module", pipeline.Module))
	if err := RunCommand(pipeline.StartCommand.ExecLine, true); err != nil {
		return exec.ErrNotFound
	}
	pipeline.Redis.PublishWithTimestamp(pipeline.StartCommand.RedisSuffix, pipeline.StartCommand.RedisMessage)

	return nil
}

// RunCommand runs or starts a command
func RunCommand(command string, start bool) error {
	c := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)

	// Log the output of the command to the console
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if start {
		if err := c.Start(); err != nil {
			return err
		}
	} else {
		if err := c.Run(); err != nil {
			return err
		}
	}

	return nil
}

// setupPushDir creates a temporary directory to do the Git operations in
func setupPushDir(srcDir, pushDir string) error {
	if err := os.RemoveAll(pushDir); err != nil {
		return err
	}

	if err := os.MkdirAll(pushDir, 0777); err != nil {
		return err
	}

	if err := copy.Copy(srcDir, pushDir); err != nil {
		return err
	}

	return nil
}
