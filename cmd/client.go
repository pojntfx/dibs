package cmd

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/pojntfx/godibs/pkg/workers"
	"github.com/spf13/cobra"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	GIT_UP_BASE_URL string

	PIPELINE_UP_DIR_SRC       string
	PIPELINE_UP_DIR_PUSH      string
	PIPELINE_UP_DIR_WATCH     string
	PIPELINE_UP_FILE_MOD      string
	PIPELINE_UP_BUILD_COMMAND string
	PIPELINE_UP_TEST_COMMAND  string
	PIPELINE_UP_START_COMMAND string
	PIPELINE_UP_REGEX_IGNORE  string

	PIPELINE_DOWN_MODULES     string
	PIPELINE_DOWN_DIR_MODULES string
)

const (
	GIT_UP_COMMIT_MESSAGE = "up_synced"
	GIT_UP_REMOTE_NAME    = "godibs-sync"
	GIT_UP_USER_NAME      = "godibs-syncer"
	GIT_UP_USER_EMAIL     = "godibs-syncer@pojtinger.space"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start the client",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the name of the module that is to be pushed
		rawGoModContent, err := ioutil.ReadFile(PIPELINE_UP_FILE_MOD)
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
		downModules := utils.GetModulesFromRawInputString(PIPELINE_DOWN_MODULES)
		// Directory to clone the local modules to
		downModulesDir := PIPELINE_DOWN_DIR_MODULES

		// Replace the modules that are specified
		moduleWithReplaces := utils.GetModuleWithReplaces(goModContent, downModules, downModulesDir)
		ioutil.WriteFile(PIPELINE_UP_FILE_MOD, []byte(moduleWithReplaces), 0777)

		// Connect to Redis
		redis := utils.Redis{
			Addr:   REDIS_URL,
			Prefix: REDIS_PREFIX,
		}
		redis.Connect()

		// Register the module
		log.Info("Registering module ...", rz.String("Module", module))
		redis.PublishWithTimestamp(REDIS_SUFFIX_UP_REGISTERED, module)

		// Unregister the module on interrupt signal
		interrupt := make(chan os.Signal, 2)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-interrupt

			log.Info("Unregistering module ...", rz.String("Module", module))
			redis.PublishWithTimestamp(REDIS_SUFFIX_UP_UNREGISTERED, module)

			log.Info("Cleaning up ...", rz.String("Module", module), rz.String("ModuleFile", PIPELINE_UP_FILE_MOD))
			rawGoModContent, err := ioutil.ReadFile(PIPELINE_UP_FILE_MOD)
			if err != nil {
				log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
			}
			goModContent := string(rawGoModContent)
			moduleWithoutReplaces := utils.GetModuleWithoutReplaces(goModContent)
			ioutil.WriteFile(PIPELINE_UP_FILE_MOD, []byte(moduleWithoutReplaces), 0777)

			os.Exit(0)
		}()

		// Setup the pipeline
		var commandStartState *exec.Cmd

		git := utils.Git{
			RemoteName:    GIT_UP_REMOTE_NAME,
			RemoteURL:     utils.GetGitURL(GIT_UP_BASE_URL, module),
			UserName:      GIT_UP_USER_NAME,
			UserEmail:     GIT_UP_USER_EMAIL,
			CommitMessage: GIT_UP_COMMIT_MESSAGE,
		}

		testCommand, buildCommand, startCommand := utils.EventedCommand{
			LogMessage:   "Running test command ...",
			ExecLine:     PIPELINE_UP_TEST_COMMAND,
			RedisSuffix:  REDIS_SUFFIX_UP_TESTED,
			RedisMessage: module,
		}, utils.EventedCommand{
			LogMessage:   "Running build command ...",
			ExecLine:     PIPELINE_UP_BUILD_COMMAND,
			RedisSuffix:  REDIS_SUFFIX_UP_BUILT,
			RedisMessage: module,
		}, utils.EventedCommand{
			LogMessage:   "Starting start command ...",
			ExecLine:     PIPELINE_UP_START_COMMAND,
			RedisSuffix:  REDIS_SUFFIX_UP_STARTED,
			RedisMessage: module,
		}

		pipeline := utils.Pipeline{
			Module:                  module,
			ModulePushedRedisSuffix: REDIS_SUFFIX_UP_PUSHED,
			SrcDir:                  PIPELINE_UP_DIR_SRC,
			PushDir:                 PIPELINE_UP_DIR_PUSH,
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
			Modules:       downModules,
			Pipeline:      pipeline,
			LocalCloneDir: PIPELINE_DOWN_DIR_MODULES,
			Redis:         redis,
			RedisSuffix:   REDIS_SUFFIX_UP_PUSHED,
			HTTPBaseURL:   GIT_UP_BASE_URL,
		}

		// Create channels
		pipelineUpdateWorkerErrors := make(chan error, 0)
		pipelineUpdateWorkerEvents := make(chan utils.Event, 0)

		// Start worker
		go pipelineUpdateWorker.Start(pipelineUpdateWorkerErrors, pipelineUpdateWorkerEvents)

		// Create a new folder watcher
		folderWatcher := utils.FolderWatcher{
			WatchDir:    PIPELINE_UP_DIR_WATCH,
			IgnoreRegex: PIPELINE_UP_REGEX_IGNORE,
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

	},
}

func init() {
	clientCmd.PersistentFlags().StringVar(&GIT_UP_BASE_URL, "git-base-url", "http://localhost:25000/repos", "Base URL of the sync remote")

	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_DIR_SRC, "dir-src", ".", "Directory in which the source code of the module to push resides")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_DIR_PUSH, "dir-push", "/tmp/.push", "Temporary directory to put the module into before pushing")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_DIR_WATCH, "dir-watch", ".", "Directory to watch for changes")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_FILE_MOD, "modules-file", "go.mod", "Go module file of the module to push")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_BUILD_COMMAND, "cmd-build", "go build ./...", "Command to run to build the module")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_TEST_COMMAND, "cmd-test", "go test ./...", "Command to run to test the module")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_START_COMMAND, "cmd-start", "go run main.go", "Command to run to start the module")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_UP_REGEX_IGNORE, "regex-ignore", "*.pb.go", "Regular expression files to ignore")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_DOWN_MODULES, "modules-pull", "", "Comma-seperated list of the names of the modules to pull")
	clientCmd.PersistentFlags().StringVar(&PIPELINE_DOWN_DIR_MODULES, "dir-pull", "/tmp/modules", "Directory to pull the modules to")

	rootCmd.AddCommand(clientCmd)
}
