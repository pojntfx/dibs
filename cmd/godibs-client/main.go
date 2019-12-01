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

	first := make(chan struct{}, 1)
	first <- struct{}{}
	var commandStartState *exec.Cmd

	err = utils.RunPipeline(r, m, commandStartState, config.REDIS_CHANNEL_PREFIX, config.REDIS_CHANNEL_MODULE_BUILT, config.REDIS_CHANNEL_MODULE_TESTED, config.REDIS_CHANNEL_MODULE_PUSHED, config.REDIS_CHANNEL_MODULE_STARTED, config.COMMAND_BUILD, config.COMMAND_TEST, config.COMMAND_START, config.GIT_BASE_URL, config.GIT_REMOTE_NAME, config.GIT_NAME, config.GIT_EMAIL, config.GIT_COMMIT_MESSAGE, config.SRC_DIR, config.PUSH_DIR)
	if err != nil {
		panic(err)
	}

	for w.IsRunning() {
		select {
		case <-first:
		case <-w.ChangeDetails():
			err := utils.RunPipeline(r, m, commandStartState, config.REDIS_CHANNEL_PREFIX, config.REDIS_CHANNEL_MODULE_BUILT, config.REDIS_CHANNEL_MODULE_TESTED, config.REDIS_CHANNEL_MODULE_PUSHED, config.REDIS_CHANNEL_MODULE_STARTED, config.COMMAND_BUILD, config.COMMAND_TEST, config.COMMAND_START, config.GIT_BASE_URL, config.GIT_REMOTE_NAME, config.GIT_NAME, config.GIT_EMAIL, config.GIT_COMMIT_MESSAGE, config.SRC_DIR, config.PUSH_DIR)

			if err != nil {
				panic(err)
			}
		}
	}
}
