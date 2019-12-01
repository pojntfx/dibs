package main

import (
	"github.com/pojntfx/godibs/pkg/config"
	"github.com/pojntfx/godibs/pkg/utils"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// TODO:
// - Use struct for pipeline and the pipeline only (only define it once, then use the `.Run()` function twice)

func main() {
	err, m := utils.GetModuleName(config.SYNC_MODULE_PUSH_MOD_FILE)
	if err != nil {
		panic(err)
	}

	redis := utils.Redis{
		Addr:   config.REDIS_URL,
		Prefix: config.REDIS_CHANNEL_PREFIX,
	}
	redis.Connect()

	log.Info("Registering module ...")
	redis.PublishWithTimestamp(config.REDIS_CHANNEL_MODULE_REGISTERED, m)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Unregistering module ...")
		redis.PublishWithTimestamp(config.REDIS_CHANNEL_MODULE_UNREGISTERED, m)
		os.Exit(0)
	}()

	w := utils.GetNewFolderWatcher(config.SYNC_MODULE_PUSH_WATCH_GLOB, config.PUSH_DIR)

	var commandStartState *exec.Cmd

	git := utils.Git{
		RemoteName:    config.GIT_REMOTE_NAME,
		RemoteURL:     utils.GetGitURL(config.GIT_BASE_URL, m),
		UserName:      config.GIT_NAME,
		UserEmail:     config.GIT_EMAIL,
		CommitMessage: config.GIT_COMMIT_MESSAGE,
	}

	testCommand := utils.EventedCommand{
		LogMessage:   "Running test command",
		ExecLine:     config.COMMAND_TEST,
		RedisSuffix:  config.REDIS_CHANNEL_MODULE_TESTED,
		RedisMessage: m,
	}

	buildCommand := utils.EventedCommand{
		LogMessage:   "Running build command",
		ExecLine:     config.COMMAND_BUILD,
		RedisSuffix:  config.REDIS_CHANNEL_MODULE_BUILT,
		RedisMessage: m,
	}

	startCommand := utils.EventedCommand{
		LogMessage:   "Starting start command",
		ExecLine:     config.COMMAND_START,
		RedisSuffix:  config.REDIS_CHANNEL_MODULE_STARTED,
		RedisMessage: m,
	}

	pipeline := utils.Pipeline{
		Module:                  m,
		ModulePushedRedisSuffix: config.REDIS_CHANNEL_MODULE_PUSHED,
		SrcDir:                  config.SRC_DIR,
		PushDir:                 config.PUSH_DIR,
		RunCommands:             []utils.EventedCommand{testCommand, buildCommand},
		StartCommand:            startCommand,
		StartCommandState:       commandStartState,
		Git:                     git,
		Redis:                   redis,
	}

	if err := pipeline.RunAll(); err != nil {
		panic(err) // TODO: Log fatal instead of panic
	}

	for w.IsRunning() {
		select {
		case <-w.ChangeDetails():
			if err := pipeline.RunAll(); err != nil {
				panic(err) // TODO: Log fatal instead of panic
			}
		}
	}
}
