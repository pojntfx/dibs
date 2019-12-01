package main

import (
	"github.com/pojntfx/godibs/pkg/config"
	"github.com/pojntfx/godibs/pkg/utils"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	// Get the name of the module that is to be pushed
	err, module := utils.GetModuleName(config.PIPELINE_UP_MOD_FILE)
	if err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err), rz.String("Module", module))
	}

	// Connect to Redis
	redis := utils.Redis{
		Addr:   config.REDIS_URL,
		Prefix: config.REDIS_PREFIX,
	}
	redis.Connect()

	// Register the module
	log.Info("Registering module ...", rz.String("Module", module))
	redis.PublishWithTimestamp(config.REDIS_SUFFIX_UP_REGISTERED, module)

	// Unregister the module on interrupt signal
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt

		log.Info("Unregistering module ...", rz.String("Module", module))
		redis.PublishWithTimestamp(config.REDIS_SUFFIX_UP_UNREGISTERED, module)

		os.Exit(0)
	}()

	// Setup the pipeline
	var commandStartState *exec.Cmd

	git := utils.Git{
		RemoteName:    config.GIT_UP_REMOTE_NAME,
		RemoteURL:     utils.GetGitURL(config.GIT_BASE_URL, module),
		UserName:      config.GIT_UP_USER_NAME,
		UserEmail:     config.GIT_UP_USER_EMAIL,
		CommitMessage: config.GIT_UP_COMMIT_MESSAGE,
	}

	testCommand, buildCommand, startCommand := utils.EventedCommand{
		LogMessage:   "Running test command ...",
		ExecLine:     config.PIPELINE_UP_TEST_COMMAND,
		RedisSuffix:  config.REDIS_SUFFIX_UP_TESTED,
		RedisMessage: module,
	}, utils.EventedCommand{
		LogMessage:   "Running build command ...",
		ExecLine:     config.PIPELINE_UP_BUILD_COMMAND,
		RedisSuffix:  config.REDIS_SUFFIX_UP_BUILT,
		RedisMessage: module,
	}, utils.EventedCommand{
		LogMessage:   "Starting start command ...",
		ExecLine:     config.PIPELINE_UP_START_COMMAND,
		RedisSuffix:  config.REDIS_SUFFIX_UP_STARTED,
		RedisMessage: module,
	}

	pipeline := utils.Pipeline{
		Module:                  module,
		ModulePushedRedisSuffix: config.REDIS_SUFFIX_UP_PUSHED,
		SrcDir:                  config.PIPELINE_UP_SRC_DIR,
		PushDir:                 config.PIPELINE_UP_PUSH_DIR,
		RunCommands:             []utils.EventedCommand{testCommand, buildCommand},
		StartCommand:            startCommand,
		StartCommandState:       commandStartState,
		Git:                     git,
		Redis:                   redis,
	}

	// Run the pipeline once. If there are errors, don't exit
	if err := pipeline.RunAll(); err != nil {
		log.Error("Error", rz.String("System", "Client"), rz.String("Module", module), rz.Err(err))
	}

	// Create a new folder watcher
	folderWatcher := utils.FolderWatcher{
		WatchDir:  config.PIPELINE_UP_WATCH_DIR,
		IgnoreDir: config.PIPELINE_UP_PUSH_DIR,
	}
	folderWatcher.Start()

	// Start the main loop
	for folderWatcher.FolderWatcher.IsRunning() {
		select {
		case <-folderWatcher.FolderWatcher.ChangeDetails():
			// Run the pipeline again on every file change. If there are errors, don't exit
			if err := pipeline.RunAll(); err != nil {
				log.Error("Error", rz.String("System", "Client"), rz.String("Module", module), rz.Err(err))
			}
		}
	}
}
