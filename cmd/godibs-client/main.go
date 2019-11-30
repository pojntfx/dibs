package main

import (
	"errors"
	fswatch "github.com/andreaskoch/go-fswatch"
	"github.com/go-redis/redis/v7"
	"github.com/otiai10/copy"
	"github.com/pojntfx/godibs/pkg/config"
	"github.com/pojntfx/godibs/pkg/utils"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	git "gopkg.in/src-d/go-git.v4"
	gitconf "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	REDIS_URL                   = os.Getenv("REDIS_URL")
	REDIS_CHANNEL_PREFIX        = os.Getenv("REDIS_CHANNEL_PREFIX")
	GIT_BASE_URL                = os.Getenv("GIT_BASE_URL")
	GIT_REMOTE_NAME             = os.Getenv("GIT_REMOTE_NAME")
	GIT_NAME                    = os.Getenv("GIT_NAME")
	GIT_EMAIL                   = os.Getenv("GIT_EMAIL")
	SRC_DIR                     = os.Getenv("SRC_DIR")
	PUSH_DIR                    = os.Getenv("PUSH_DIR")
	SYNC_MODULE_PUSH_MOD_FILE   = os.Getenv("SYNC_MODULE_PUSH_MOD_FILE")
	SYNC_MODULE_PUSH_WATCH_GLOB = os.Getenv("SYNC_MODULE_PUSH_WATCH_GLOB")
	COMMAND_BUILD               = os.Getenv("COMMAND_BUILD")
	COMMAND_TEST                = os.Getenv("COMMAND_TEST")
	COMMAND_START               = os.Getenv("COMMAND_START")
)

func main() {
	err, m := GetModuleName(SYNC_MODULE_PUSH_MOD_FILE)
	if err != nil {
		panic(err)
	}

	r := utils.NewRedisClient(REDIS_URL)

	log.Info("Registering module ...")
	RegisterModule(r, REDIS_CHANNEL_PREFIX, m)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Unregistering module ...")
		UnregisterModule(r, REDIS_CHANNEL_PREFIX, m)
		os.Exit(0)
	}()

	w := GetNewFolderWatcher(SYNC_MODULE_PUSH_WATCH_GLOB, PUSH_DIR)

	first := make(chan struct{}, 1)
	first <- struct{}{}
	var commandStartState *exec.Cmd

	err = RunPipeline(r, m, commandStartState, REDIS_CHANNEL_PREFIX, config.REDIS_CHANNEL_MODULE_BUILT, config.REDIS_CHANNEL_MODULE_TESTED, config.REDIS_CHANNEL_MODULE_PUSHED, config.REDIS_CHANNEL_MODULE_STARTED, COMMAND_BUILD, COMMAND_TEST, COMMAND_START, GIT_BASE_URL, GIT_REMOTE_NAME, GIT_NAME, GIT_EMAIL, SRC_DIR, PUSH_DIR)
	if err != nil {
		panic(err)
	}

	for w.IsRunning() {
		select {
		case <-first:
		case <-w.ChangeDetails():
			err := RunPipeline(r, m, commandStartState, REDIS_CHANNEL_PREFIX, config.REDIS_CHANNEL_MODULE_BUILT, config.REDIS_CHANNEL_MODULE_TESTED, config.REDIS_CHANNEL_MODULE_PUSHED, config.REDIS_CHANNEL_MODULE_STARTED, COMMAND_BUILD, COMMAND_TEST, COMMAND_START, GIT_BASE_URL, GIT_REMOTE_NAME, GIT_NAME, GIT_EMAIL, SRC_DIR, PUSH_DIR)
			if err != nil {
				panic(err)
			}
		}
	}
}

// RunPipeline runs the entire pipeline
func RunPipeline(r *redis.Client, m string, commandStartState *exec.Cmd, channelPrefix, moduleBuildSuffix, moduleTestedSuffix, modulePushedSuffix, moduleStartedSuffix, commandBuild, commandTest, commandStart, gitBaseUrl, gitRemoteName, gitName, gitEmail, srcDir, pushDir string) error {
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
	err = PushModule(r, channelPrefix, modulePushedSuffix, m, pushDir, gitRemoteName, pushURL, gitName, gitEmail)
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

// GetGitURL returns the URL of a git repo
func GetGitURL(baseURL, m string) string {
	return baseURL + "/" + m
}

// GetModuleName returns the module name from `go.mod`
func GetModuleName(goModFilePath string) (error, string) {
	f, err := ioutil.ReadFile(goModFilePath)
	if err != nil {
		return errors.New("Could not read module file"), ""
	}

	for _, line := range strings.Split(string(f), "\n") {
		if strings.Contains(line, "module") {
			return nil, strings.Split(line, "module ")[1]
		}
	}

	return errors.New("Could find module declaration"), ""
}

// GetNewFolderWatcher returns a new folder watcher
func GetNewFolderWatcher(watchGlob, pushDir string) *fswatch.FolderWatcher {
	w := fswatch.NewFolderWatcher(watchGlob, true, func(path string) bool { return strings.Contains(path, pushDir) }, 1)
	w.Start()

	return w
}

// RegisterModule registers a module in Redis
func RegisterModule(r *redis.Client, prefix, m string) {
	r.Publish(prefix+":"+config.REDIS_CHANNEL_MODULE_REGISTERED, withTimestamp(m))
}

// UnregisterModule unregisters a module from Redis
func UnregisterModule(r *redis.Client, prefix, m string) {
	r.Publish(prefix+":"+config.REDIS_CHANNEL_MODULE_UNREGISTERED, withTimestamp(m))
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
	r.Publish(prefix+":"+suffix, withTimestamp(m))
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

// PushModule adds all files to a git repo, commits and finally pushes them to a remote
func PushModule(r *redis.Client, prefix, suffix, m, pushDir, gitRemoteName, pushURL, gitName, gitEmail string) error {
	g, err := git.PlainOpen(filepath.Join(pushDir))
	if err != nil {
		return err
	}

	_, err = g.CreateRemote(&gitconf.RemoteConfig{
		Name: gitRemoteName,
		URLs: []string{pushURL},
	})
	if err != nil {
		return err
	}

	wt, err := g.Worktree()
	if err != nil {
		return err
	}

	_, err = wt.Add(".")
	if err != nil {
		return err
	}

	_, err = wt.Commit(withTimestamp(config.GIT_COMMIT_MESSAGE), &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitName,
			Email: gitEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	err = g.Push(&git.PushOptions{
		RemoteName: gitRemoteName,
		RefSpecs:   []gitconf.RefSpec{"+refs/heads/master:refs/heads/master"},
	})
	if err != nil {
		return err
	}

	r.Publish(prefix+":"+suffix, withTimestamp(m))

	return nil
}

// withTimestamp gets a message name with the current timestamp
func withTimestamp(m string) string {
	t := time.Now().UnixNano()
	return m + "@" + strconv.Itoa(int(t))
}
