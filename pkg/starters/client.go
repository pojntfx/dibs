package starters

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/pojntfx/dibs/pkg/workers"
	"github.com/radovskyb/watcher"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Client is a client for the sync server
type Client struct {
	PipelineUpFileMod                string // Go module file of the module to push
	PipelineDownModules              string // Comma-separated list of the names of the modules to pull
	PipelineDownDirModules           string // Directory to pull the modules to
	PipelineUpUnitTestCommand        string // Command to run to unit test the module
	PipelineUpIntegrationTestCommand string // Command to run to integration test the module
	PipelineUpBuildCommand           string // Command to run to build the module
	PipelineUpStartCommand           string // Command to run to start the module
	PipelineUpDirSrc                 string // Directory in which the source code of the module to push resides
	PipelineUpDirPush                string // Temporary directory to put the module into before pushing
	PipelineUpDirWatch               string // Directory to watch for changes
	PipelineUpRegexIgnore            string // Regular expression for files to ignore

	RedisUrl                  string // URL of the Redis instance to use
	RedisPrefix               string // Redis channel prefix
	RedisPassword             string // Redis password
	RedisSuffixUpRegistered   string // Redis channel suffix for "module registered" messages
	RedisSuffixUpUnRegistered string // Redis channel suffix for "module unregistered" messages
	RedisSuffixUpTested       string // Redis channel suffix for "module tested" messages
	RedisSuffixUpBuilt        string // Redis channel suffix for "module built" messages
	RedisSuffixUpStarted      string // Redis channel suffix for "module started" messages
	RedisSuffixUpPushed       string // Redis channel suffix for "module pushed" messages

	GitUpRemoteName    string // Name of the sync remote to add
	GitUpBaseURL       string // Base URL of the sync remote
	GitUpUserName      string // Username for Git commits
	GitUpUserEmail     string // Email for Git commits
	GitUpCommitMessage string // Message for Git commits
}

// Start starts the sync client
func (client *Client) Start() {
	// Get the name of the module that is to be pushed
	rawGoModContent, err := ioutil.ReadFile(client.PipelineUpFileMod)
	if err != nil {
		utils.LogErrorFatal("Error", err)
	}
	// Get the content of the Go module file
	goModContent := string(rawGoModContent)
	err, module := utils.GetModuleName(goModContent)
	if err != nil {
		utils.LogErrorForModuleFatal("Error", err, module)
	}
	// Get the modules that are to be downloaded
	downModules := utils.GetModulesFromRawInputString(client.PipelineDownModules)
	// Directory to clone the local modules to
	downModulesDir := client.PipelineDownDirModules

	// Replace the modules that are specified. Don't run if no pull modules have been specified.
	if downModules[0] != "" {
		moduleWithReplaces, err := utils.GetModuleWithReplaces(goModContent, downModules, downModulesDir)
		if err != nil {
			utils.LogErrorForModuleFatal("Error", err, module)
		}
		if err := ioutil.WriteFile(client.PipelineUpFileMod, []byte(moduleWithReplaces), 0777); err != nil {
			utils.LogErrorForModuleFatal("Error", err, module)
		}
	}

	// Connect to Redis
	redis := utils.Redis{
		Addr:     client.RedisUrl,
		Prefix:   client.RedisPrefix,
		Password: client.RedisPassword,
	}
	redis.Connect()

	// Register the module
	utils.LogForModule("Registering module", module)
	redis.PublishWithTimestamp(client.RedisSuffixUpRegistered, module)

	// Setup the pipeline
	var commandStartState *exec.Cmd

	git := utils.Git{
		RemoteName:    client.GitUpRemoteName,
		RemoteURL:     utils.GetGitURL(client.GitUpBaseURL, module),
		UserName:      client.GitUpUserName,
		UserEmail:     client.GitUpUserEmail,
		CommitMessage: client.GitUpCommitMessage,
	}

	unitTestCommand, integrationTestCommand, buildCommand, startCommand := utils.CommandWithEvent{
		LogMessage:   "Unit testing module",
		ExecLine:     client.PipelineUpUnitTestCommand,
		RedisSuffix:  client.RedisSuffixUpTested,
		RedisMessage: module,
	}, utils.CommandWithEvent{
		LogMessage:   "Integration testing module",
		ExecLine:     client.PipelineUpIntegrationTestCommand,
		RedisSuffix:  client.RedisSuffixUpTested,
		RedisMessage: module,
	}, utils.CommandWithEvent{
		LogMessage:   "Building module",
		ExecLine:     client.PipelineUpBuildCommand,
		RedisSuffix:  client.RedisSuffixUpBuilt,
		RedisMessage: module,
	}, utils.CommandWithEvent{
		LogMessage:   "Starting module",
		ExecLine:     client.PipelineUpStartCommand,
		RedisSuffix:  client.RedisSuffixUpStarted,
		RedisMessage: module,
	}

	pipeline := utils.Pipeline{
		Module:                  module,
		ModulePushedRedisSuffix: client.RedisSuffixUpPushed,
		SrcDir:                  client.PipelineUpDirSrc,
		PushDir:                 client.PipelineUpDirPush,
		RunCommands:             []utils.CommandWithEvent{unitTestCommand, integrationTestCommand, buildCommand},
		StartCommand:            startCommand,
		StartCommandState:       commandStartState,
		Git:                     git,
		Redis:                   redis,
	}

	// Unregister the module on interrupt signal
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt

		utils.LogForModule("Stopping module", module)
		processGroupId, err := syscall.Getpgid(pipeline.StartCommandState.Process.Pid)
		if err != nil {
			utils.LogErrorFatalCouldStopModule(err)
		}

		if err := syscall.Kill(-processGroupId, syscall.SIGKILL); err != nil {
			utils.LogErrorFatalCouldStopModule(err)
		}

		utils.LogForModule("Unregistering module", module)
		redis.PublishWithTimestamp(client.RedisSuffixUpUnRegistered, module)

		utils.LogForModule("Cleaning up", module)

		// Remove the added replace directives. Don't run if no pull modules have been specified.
		if downModules[0] != "" {
			rawGoModContent, err := ioutil.ReadFile(client.PipelineUpFileMod)
			if err != nil {
				utils.LogErrorFatal("Error", err)
			}
			goModContent := string(rawGoModContent)
			moduleWithoutReplaces := utils.GetModuleWithoutReplaces(goModContent)
			if err := ioutil.WriteFile(client.PipelineUpFileMod, []byte(moduleWithoutReplaces), 0777); err != nil {
				utils.LogErrorFatal("Error", err)
			}
		}

		os.Exit(0)
	}()

	// Run the pipeline once. If there are errors, don't exit.
	if err := pipeline.RunAll(); err != nil {
		utils.LogErrorForModuleFatal("Error", err, module)
	}

	// Setup worker
	pipelineUpdateWorker := &workers.PipelineUpdateWorker{
		Modules:       downModules,
		Pipeline:      pipeline,
		LocalCloneDir: client.PipelineDownDirModules,
		Redis:         redis,
		RedisSuffix:   client.RedisSuffixUpPushed,
		HTTPBaseURL:   client.GitUpBaseURL,
	}

	// Create channels
	pipelineUpdateWorkerErrors := make(chan error, 0)
	pipelineUpdateWorkerEvents := make(chan utils.Event, 0)

	// Start worker. Don't run if no pull modules have been specified.
	if downModules[0] != "" {
		go pipelineUpdateWorker.Start(pipelineUpdateWorkerErrors, pipelineUpdateWorkerEvents)
	}

	// Create a new folder watcher
	folderWatcher := utils.FolderWatcher{
		WatchDir:    client.PipelineUpDirWatch,
		IgnoreRegex: client.PipelineUpRegexIgnore,
	}

	// Register the folder watcher's event handlers
	if err := folderWatcher.Start(func(err error) {
		utils.LogErrorFatal("Error", err)
	}, func(event watcher.Event) {
		// Run the pipeline again on every file change. If there are errors, don't exit.
		if err := pipeline.RunAll(); err != nil {
			utils.LogErrorForModule("Error", err, module)
		}
	}); err != nil {
		utils.LogErrorFatal("Error", err)
	}

	// Start the main loop
	for {
		select {
		// If there are errors, log the errors and exit
		case err := <-pipelineUpdateWorkerErrors:
			log.Fatal("Error", rz.String("system", "PipelineUpdateWorker"), rz.Err(err))
		case event := <-pipelineUpdateWorkerEvents:
			switch event.Code {
			case 0:
				log.Info("Started", rz.String("system", "PipelineUpdateWorker"), rz.String("eventMessage", event.Message))
			case 1:
				log.Info("Request", rz.String("system", "PipelineUpdateWorker"), rz.String("eventMessage", event.Message))
			case 2:
				log.Info("Stopped", rz.String("system", "PipelineUpdateWorker"), rz.String("eventMessage", event.Message))
				return
			default:
				log.Fatal("Unknown event code", rz.String("system", "GitHTTPWorker"), rz.Int("eventCode", event.Code), rz.String("statusMessage", event.Message))
			}
		}
	}
}
