package utils

import (
	fswatch "github.com/andreaskoch/go-fswatch"
	redis "github.com/go-redis/redis/v7"
	"github.com/otiai10/copy"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// RunPipeline runs the entire pipeline
func RunPipeline(r *redis.Client, m string, commandStartState *exec.Cmd, channelPrefix, moduleBuildSuffix, moduleTestedSuffix, modulePushedSuffix, moduleStartedSuffix, commandBuild, commandTest, commandStart, gitBaseUrl, gitRemoteName, gitName, gitEmail, gitCommitMessage, srcDir, pushDir string) error {
	log.Info("Stopping module ...")
	if commandStartState != nil {
		commandStartState.Process.Kill()
	}

	log.Info("Copying module ...")
	err := SetupPushDir(srcDir, pushDir)
	if err != nil {
		return err
	}

	log.Info("Building module ...")
	err = RunCommand(r, channelPrefix, moduleTestedSuffix, m, commandBuild, false)
	if err != nil {
		return err
	}

	log.Info("Testing module ...")
	err = RunCommand(r, channelPrefix, moduleTestedSuffix, m, commandTest, false)
	if err != nil {
		return err
	}

	log.Info("Pushing module ...")
	pushURL := GetGitURL(gitBaseUrl, m)
	git := Git{
		RemoteName:    gitRemoteName,
		RemoteURL:     pushURL,
		UserName:      gitName,
		UserEmail:     gitEmail,
		CommitMessage: gitCommitMessage,
	}
	err = git.PushModule(r, channelPrefix, modulePushedSuffix, m, pushDir)
	if err != nil {
		return err
	}

	log.Info("Starting module ...")
	err = RunCommand(r, channelPrefix, moduleStartedSuffix, m, commandStart, true)
	if err != nil {
		return err
	}

	return nil
}

// GetNewFolderWatcher returns a new folder watcher
func GetNewFolderWatcher(watchGlob, pushDir string) *fswatch.FolderWatcher {
	w := fswatch.NewFolderWatcher(watchGlob, true, func(path string) bool { return strings.Contains(path, pushDir) }, 1)
	w.Start()

	return w
}

// RunCommand runs or starts a command creates a corresponding message in Redis
func RunCommand(r *redis.Client, prefix, suffix, m, command string, start bool) error {
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
	r.Publish(prefix+":"+suffix, WithTimestamp(m))
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
