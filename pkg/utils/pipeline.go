package utils

import (
	"github.com/otiai10/copy"
	"os"
	"os/exec"
	"strings"
)

// CommandWithEvent is a command that also publishes an event
type CommandWithEvent struct {
	LogMessage   string // The message to log to stdout
	ExecLine     string // The command to run/start
	RedisSuffix  string // The Redis channel suffix to use for the event
	RedisMessage string // The Redis message to send
}

// Pipeline is a development configuration
type Pipeline struct {
	Module                  string             // The module that is being pushed
	ModulePushedRedisSuffix string             // The Redis suffix channel suffix to use for the pushed event
	SrcDir                  string             // The directory of the module's source
	PushDir                 string             // The temporary directory to use for the modification of the Git repo
	RunCommands             []CommandWithEvent // The commands to run
	StartCommand            CommandWithEvent   // The command to start (will keep running, but can be killed)
	StartCommandState       *exec.Cmd          // Stores the state of the start command (to make it possible to kill it)
	Git                     Git                // Git instance to use to modify the Git repo in PushDir
	Redis                   Redis              // Redis instance to use to publish the events
}

// RunCommandsOnly only runs the commands
func (pipeline *Pipeline) RunCommandsOnly() error {
	for _, runCommand := range pipeline.RunCommands {
		LogForModule(runCommand.LogMessage, pipeline.Module)
		if err := pipeline.runCommand(runCommand.ExecLine, false); err != nil {
			return err
		}
		pipeline.Redis.PublishWithTimestamp(runCommand.RedisSuffix, runCommand.RedisMessage)
	}

	return nil
}

// RunAll runs the entire pipeline
func (pipeline *Pipeline) RunAll() error {

	if pipeline.StartCommandState != nil {
		if pipeline.StartCommandState.Process != nil {
			LogForModule("Restarting pipeline", pipeline.Module)

			LogForModule("Stopping module", pipeline.Module)

			if err := pipeline.StartCommandState.Process.Kill(); err != nil {
				LogError("Could not stop module", err)
			}
		}
	} else {
		LogForModule("Starting pipeline", pipeline.Module)
	}

	LogForModule("Pushing module", pipeline.Module)

	if err := setupPushDir(pipeline.SrcDir, pipeline.PushDir); err != nil {
		return err
	}

	if err := pipeline.Git.PushToRemote(pipeline.PushDir); err != nil {
		return err
	}
	pipeline.Redis.PublishWithTimestamp(pipeline.ModulePushedRedisSuffix, pipeline.Module)

	if err := pipeline.RunCommandsOnly(); err != nil {
		LogError("Could not run pipeline's command", err)
	}

	LogForModule(pipeline.StartCommand.LogMessage, pipeline.Module)
	if err := pipeline.runCommand(pipeline.StartCommand.ExecLine, true); err != nil {
		return exec.ErrNotFound
	}
	pipeline.Redis.PublishWithTimestamp(pipeline.StartCommand.RedisSuffix, pipeline.StartCommand.RedisMessage)

	return nil
}

// runCommand runs or starts a command
func (pipeline *Pipeline) runCommand(command string, start bool) error {
	c := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)

	// Log the output of the command to the console
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if start {
		pipeline.StartCommandState = c

		return c.Start()
	} else {
		return c.Run()
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
