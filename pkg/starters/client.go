package starters

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/pojntfx/dibs/pkg/workers"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Client is a client for dibs server
type Client struct {
	PipelineUpFileMod      string // Go module file of the module to push
	PipelineDownModules    string // Comma-separated list of the names of the modules to pull
	PipelineDownDirModules string // Directory to pull the modules to
	PipelineUpTestCommand  string // Command to run to test the module
	PipelineUpBuildCommand string // Command to run to build the module
	PipelineUpStartCommand string // Command to run to start the module
	PipelineUpDirSrc       string // Directory in which the source code of the module to push resides
	PipelineUpDirPush      string // Temporary directory to put the module into before pushing
	PipelineUpDirWatch     string // Directory to watch for changes
	PipelineUpRegexIgnore  string // Regular expression for files to ignore

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

// Start starts the client
func (client *Client) Start() {
	// Get the name of the module that is to be pushed
	rawGoModContent, err := ioutil.ReadFile(client.PipelineUpFileMod)
	if err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
	}
	// Get the content of the Go module file
	goModContent := string(rawGoModContent)
	err, module := utils.GetModuleName(goModContent)
	if err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err), rz.String("Module", module))
	}
	// Get the modules that are to be downloaded
	downModules := utils.GetModulesFromRawInputString(client.PipelineDownModules)
	// Directory to clone the local modules to
	downModulesDir := client.PipelineDownDirModules

	// Replace the modules that are specified. Don't run if no pull modules have been specified.
	if downModules[0] != "" {
		moduleWithReplaces, err := utils.GetModuleWithReplaces(goModContent, downModules, downModulesDir)
		if err != nil {
			log.Fatal("Error", rz.String("System", "Client"), rz.Err(err), rz.String("Module", module))
		}
		if err := ioutil.WriteFile(client.PipelineUpFileMod, []byte(moduleWithReplaces), 0777); err != nil {
			log.Fatal("Error", rz.String("System", "Client"), rz.Err(err), rz.String("Module", module))
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
	log.Info("Registering module ...", rz.String("Module", module))
	redis.PublishWithTimestamp(client.RedisSuffixUpRegistered, module)

	// Unregister the module on interrupt signal
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt

		log.Info("Unregistering module ...", rz.String("Module", module))
		redis.PublishWithTimestamp(client.RedisSuffixUpUnRegistered, module)

		log.Info("Cleaning up ...", rz.String("Module", module), rz.String("ModuleFile", client.PipelineUpFileMod))

		// Remove the added replace directives. Don't run if no pull modules have been specified.
		if downModules[0] != "" {
			rawGoModContent, err := ioutil.ReadFile(client.PipelineUpFileMod)
			if err != nil {
				log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
			}
			goModContent := string(rawGoModContent)
			moduleWithoutReplaces := utils.GetModuleWithoutReplaces(goModContent)
			if err := ioutil.WriteFile(client.PipelineUpFileMod, []byte(moduleWithoutReplaces), 0777); err != nil {
				log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
			}
		}

		os.Exit(0)
	}()

	// Setup the pipeline
	var commandStartState *exec.Cmd

	git := utils.Git{
		RemoteName:    client.GitUpRemoteName,
		RemoteURL:     utils.GetGitURL(client.GitUpBaseURL, module),
		UserName:      client.GitUpUserName,
		UserEmail:     client.GitUpUserEmail,
		CommitMessage: client.GitUpCommitMessage,
	}

	testCommand, buildCommand, startCommand := utils.CommandWithEvent{
		LogMessage:   "Running test command ...",
		ExecLine:     client.PipelineUpTestCommand,
		RedisSuffix:  client.RedisSuffixUpTested,
		RedisMessage: module,
	}, utils.CommandWithEvent{
		LogMessage:   "Running build command ...",
		ExecLine:     client.PipelineUpBuildCommand,
		RedisSuffix:  client.RedisSuffixUpBuilt,
		RedisMessage: module,
	}, utils.CommandWithEvent{
		LogMessage:   "Starting start command ...",
		ExecLine:     client.PipelineUpStartCommand,
		RedisSuffix:  client.RedisSuffixUpStarted,
		RedisMessage: module,
	}

	pipeline := utils.Pipeline{
		Module:                  module,
		ModulePushedRedisSuffix: client.RedisSuffixUpPushed,
		SrcDir:                  client.PipelineUpDirSrc,
		PushDir:                 client.PipelineUpDirPush,
		RunCommands:             []utils.CommandWithEvent{testCommand, buildCommand},
		StartCommand:            startCommand,
		StartCommandState:       commandStartState,
		Git:                     git,
		Redis:                   redis,
	}

	// Run the pipeline once. If there are errors, don't exit.
	if err := pipeline.RunAll(); err != nil {
		log.Error("Error", rz.String("System", "Client"), rz.String("Module", module), rz.Err(err))
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
	folderWatcher.Start()

	// Start the main loop
	for folderWatcher.FolderWatcher.IsRunning() {
		select {
		// If there are errors, log the errors and exit
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
			// Run the pipeline again on every file change. If there are errors, don't exit.
			if err := pipeline.RunAll(); err != nil {
				log.Error("Error", rz.String("System", "Client"), rz.String("Module", module), rz.Err(err))
			}
		}
	}
}
