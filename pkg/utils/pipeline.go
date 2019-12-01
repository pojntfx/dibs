package utils

import (
	fswatch "github.com/andreaskoch/go-fswatch"
	"github.com/otiai10/copy"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type EventedCommand struct {
	LogMessage   string
	ExecLine     string
	RedisSuffix  string
	RedisMessage string
}

type Pipeline struct {
	Module                  string
	ModulePushedRedisSuffix string
	SrcDir                  string
	PushDir                 string
	RunCommands             []EventedCommand
	StartCommand            EventedCommand
	StartCommandState       *exec.Cmd
	Git                     Git
	Redis                   Redis
}

func (pipeline *Pipeline) RunAll() error {
	if pipeline.StartCommandState != nil {
		log.Info("Stopping module ...")
		pipeline.StartCommandState.Process.Kill()
	}

	if err := SetupPushDir(pipeline.SrcDir, pipeline.PushDir); err != nil {
		return err
	}

	if err := pipeline.Git.PushModule(pipeline.Module, pipeline.PushDir); err != nil {
		return err
	}
	pipeline.Redis.PublishWithTimestamp(pipeline.ModulePushedRedisSuffix, pipeline.Module)

	for _, runCommand := range pipeline.RunCommands {
		log.Info(runCommand.LogMessage)
		if err := RunCommand(runCommand.ExecLine, false); err != nil {
			return err
		}
		pipeline.Redis.PublishWithTimestamp(runCommand.RedisSuffix, runCommand.RedisMessage)
	}

	log.Info(pipeline.StartCommand.LogMessage)
	if err := RunCommand(pipeline.StartCommand.ExecLine, true); err != nil {
		return exec.ErrNotFound
	}
	pipeline.Redis.PublishWithTimestamp(pipeline.StartCommand.RedisSuffix, pipeline.StartCommand.RedisMessage)

	return nil
}

// GetNewFolderWatcher returns a new folder watcher
func GetNewFolderWatcher(watchGlob, pushDir string) *fswatch.FolderWatcher {
	w := fswatch.NewFolderWatcher(watchGlob, true, func(path string) bool { return strings.Contains(path, pushDir) }, 1)
	w.Start()

	return w
}

// RunCommand runs or starts a command creates a corresponding message in Redis
func RunCommand(command string, start bool) error {
	c := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	var err error
	if start {
		err = c.Start()
	} else {
		err = c.Run()
	}
	if err != nil {
		return err
	}
	return nil
}

// SetupPushDir creates a temporary directory to do the git operations in
func SetupPushDir(srcDir, pushDir string) error {
	err := os.RemoveAll(pushDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(pushDir, 0777)
	if err != nil {
		return err
	}

	log.Info("Copying internal", rz.String("src", srcDir), rz.String("dist", pushDir))
	err = copy.Copy(srcDir, pushDir)
	if err != nil {
		return err
	}

	return nil
}

// WithTimestamp gets a message name with the current timestamp
func WithTimestamp(m string) string {
	t := time.Now().UnixNano()
	return m + "@" + strconv.Itoa(int(t))
}
