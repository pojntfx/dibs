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
	err, module := utils.GetModuleName(config.SYNC_MODULE_PUSH_MOD_FILE)
	if err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err), rz.String("Module", module))
	}

	// Connect to Redis
	redis := utils.Redis{
		Addr:   config.REDIS_URL,
		Prefix: config.REDIS_CHANNEL_PREFIX,
	}
	redis.Connect()

	// Register the module
	log.Info("Registering module ...", rz.String("Module", module))
	redis.PublishWithTimestamp(config.REDIS_CHANNEL_MODULE_REGISTERED, module)

	// Unregister the module on interrupt signal
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt

		log.Info("Unregistering module ...", rz.String("Module", module))
		redis.PublishWithTimestamp(config.REDIS_CHANNEL_MODULE_UNREGISTERED, module)

		os.Exit(0)
	}()

	// Setup the pipeline
	var commandStartState *exec.Cmd

	git := utils.Git{
		RemoteName:    config.GIT_REMOTE_NAME,
		RemoteURL:     utils.GetGitURL(config.GIT_BASE_URL, module),
		UserName:      config.GIT_NAME,
		UserEmail:     config.GIT_EMAIL,
		CommitMessage: config.GIT_COMMIT_MESSAGE,
	}

	testCommand, buildCommand, startCommand := utils.EventedCommand{
		LogMessage:   "Running test command ...",
		ExecLine:     config.COMMAND_TEST,
		RedisSuffix:  config.REDIS_CHANNEL_MODULE_TESTED,
		RedisMessage: module,
	}, utils.EventedCommand{
		LogMessage:   "Running build command ...",
		ExecLine:     config.COMMAND_BUILD,
		RedisSuffix:  config.REDIS_CHANNEL_MODULE_BUILT,
		RedisMessage: module,
	}, utils.EventedCommand{
		LogMessage:   "Starting start command ...",
		ExecLine:     config.COMMAND_START,
		RedisSuffix:  config.REDIS_CHANNEL_MODULE_STARTED,
		RedisMessage: module,
	}

	pipeline := utils.Pipeline{
		Module:                  module,
		ModulePushedRedisSuffix: config.REDIS_CHANNEL_MODULE_PUSHED,
		SrcDir:                  config.SRC_DIR,
		PushDir:                 config.PUSH_DIR,
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
		WatchGlob: config.SYNC_MODULE_PUSH_WATCH_GLOB,
		IgnoreDir: config.PUSH_DIR,
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
