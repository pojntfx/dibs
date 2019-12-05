package main

import (
	"github.com/pojntfx/godibs/pkg/config"
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/pojntfx/godibs/pkg/workers"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	// Get the name of the module that is to be pushed
	rawGoModContent, err := ioutil.ReadFile(config.PIPELINE_UP_FILE_MOD)
	if err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
	}
	goModContent := string(rawGoModContent)
	err, module := utils.GetModuleName(goModContent)
	if err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err), rz.String("Module", module))
	}

	// Replace the modules that are specified
	moduleWithReplaces := utils.GetModuleWithReplaces(goModContent, []string{"github.com/andreaskoch/go-fswatch"}, "localhost.localdomain:5000")
	ioutil.WriteFile(config.PIPELINE_UP_FILE_MOD, []byte(moduleWithReplaces), 0777)

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

		log.Info("Cleaning up ...", rz.String("Module", module), rz.String("ModuleFile", config.PIPELINE_UP_FILE_MOD))
		rawGoModContent, err := ioutil.ReadFile(config.PIPELINE_UP_FILE_MOD)
		if err != nil {
			log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
		}
		goModContent := string(rawGoModContent)
		moduleWithoutReplaces := utils.GetModuleWithoutReplaces(goModContent)
		ioutil.WriteFile(config.PIPELINE_UP_FILE_MOD, []byte(moduleWithoutReplaces), 0777)

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
		SrcDir:                  config.PIPELINE_UP_DIR_SRC,
		PushDir:                 config.PIPELINE_UP_DIR_PUSH,
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

	// Setup worker
	pipelineUpdateWorker := &workers.PipelineUpdateWorker{
		Pipeline:    pipeline,
		Redis:       redis,
		RedisSuffix: config.REDIS_SUFFIX_UP_PUSHED,
	}

	// Create channels
	pipelineUpdateWorkerErrors := make(chan error, 0)
	pipelineUpdateWorkerEvents := make(chan utils.Event, 0)

	// Start worker
	go pipelineUpdateWorker.Start(pipelineUpdateWorkerErrors, pipelineUpdateWorkerEvents)

	// Create a new folder watcher
	folderWatcher := utils.FolderWatcher{
		WatchDir:  config.PIPELINE_UP_DIR_WATCH,
		IgnoreDir: config.PIPELINE_UP_DIR_PUSH,
	}
	folderWatcher.Start()

	// Start the main loop
	for folderWatcher.FolderWatcher.IsRunning() {
		select {
		// If there are errors, log the erros and exit
		case err := <-pipelineUpdateWorkerErrors:
			log.Fatal("Error", rz.String("System", "PipelineUpdateWorker"), rz.Err(err))
		case event := <-pipelineUpdateWorkerEvents:
			switch event.Code {
			case 0:
				log.Info("Started", rz.String("System", "PipelineUpdateWorker"), rz.String("EventMessage", event.Message))
			case 1:
				log.Info("Request", rz.String("System", "PipelineUpdateWorker"), rz.String("EventMessage", event.Message))
			case 2:
				log.Info("Stopped", rz.String("System", "PipelineUpdateWorker"), rz.String("EventMessage", event.Message))
				return
			default:
				log.Fatal("Unknown event code", rz.String("System", "GitHTTPWorker"), rz.Int("EventCode", event.Code), rz.String("StatusMessage", event.Message))
			}
		case <-folderWatcher.FolderWatcher.ChangeDetails():
			// Run the pipeline again on every file change. If there are errors, don't exit
			if err := pipeline.RunAll(); err != nil {
				log.Error("Error", rz.String("System", "Client"), rz.String("Module", module), rz.Err(err))
			}
		}
	}
}
