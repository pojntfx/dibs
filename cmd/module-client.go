package cmd

import (
	"github.com/google/uuid"
	"github.com/pojntfx/godibs/pkg/starters"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	GitUpBaseUrl string

	PipelineUpDirSrc       string
	PipelineUpDirPush      string
	PipelineUpDirWatch     string
	PipelineUpFileMod      string
	PipelineUpBuildCommand string
	PipelineUpTestCommand  string
	PipelineUpStartCommand string
	PipelineUpRegexIgnore  string

	PipelineDownModules    string
	PipelineDownDirModules string
)

const (
	GitUpCommitMessage = "up_synced"
	GitUpRemoteName    = "godibs-sync"
	GitUpUserName      = "godibs-syncer"
	GitUpUserEmail     = "godibs-syncer@pojtinger.space"
)

// moduleClientCmd ist the command to start the client
var moduleClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start the module development client",
	Run: func(cmd *cobra.Command, args []string) {
		client := starters.Client{
			PipelineUpFileMod:      PipelineUpFileMod,
			PipelineDownModules:    PipelineDownModules,
			PipelineDownDirModules: PipelineDownDirModules,
			PipelineUpBuildCommand: PipelineUpBuildCommand,
			PipelineUpStartCommand: PipelineUpStartCommand,
			PipelineUpTestCommand:  PipelineUpTestCommand,
			PipelineUpDirSrc:       PipelineUpDirSrc,
			PipelineUpDirPush:      PipelineUpDirPush,
			PipelineUpDirWatch:     PipelineUpDirWatch,
			PipelineUpRegexIgnore:  PipelineUpRegexIgnore,

			RedisUrl:                  RedisUrl,
			RedisPrefix:               RedisPrefix,
			RedisSuffixUpRegistered:   RedisSuffixUpRegistered,
			RedisSuffixUpUnRegistered: RedisSuffixUpUnregistered,
			RedisSuffixUpTested:       RedisSuffixUpTested,
			RedisSuffixUpBuilt:        RedisSuffixUpBuilt,
			RedisSuffixUpStarted:      RedisSuffixUpStarted,
			RedisSuffixUpPushed:       RedisSuffixUpPushed,

			GitUpRemoteName:    GitUpRemoteName,
			GitUpBaseURL:       GitUpBaseUrl,
			GitUpUserName:      GitUpUserName,
			GitUpUserEmail:     GitUpUserEmail,
			GitUpCommitMessage: GitUpCommitMessage,
		}

		client.Start()
	},
}

// init maps the flags to the config
func init() {
	id := uuid.New().String()

	moduleClientCmd.PersistentFlags().StringVar(&GitUpBaseUrl, "git-base-url", "http://localhost:25000/repos", "Base URL of the sync remote")

	moduleClientCmd.PersistentFlags().StringVar(&PipelineUpDirSrc, "dir-src", ".", "Directory in which the source code of the module to push resides")
	moduleClientCmd.PersistentFlags().StringVar(&PipelineUpDirPush, "dir-push", filepath.Join(os.TempDir(), "godibs", "push", id), "Temporary directory to put the module into before pushing")
	moduleClientCmd.PersistentFlags().StringVar(&PipelineUpDirWatch, "dir-watch", ".", "Directory to watch for changes")
	moduleClientCmd.PersistentFlags().StringVar(&PipelineUpFileMod, "modules-file", "go.mod", "Go module file of the module to push")
	moduleClientCmd.PersistentFlags().StringVar(&PipelineUpBuildCommand, "cmd-build", "go build ./...", "Command to run to build the module")
	moduleClientCmd.PersistentFlags().StringVar(&PipelineUpTestCommand, "cmd-test", "go test ./...", "Command to run to test the module")
	moduleClientCmd.PersistentFlags().StringVar(&PipelineUpStartCommand, "cmd-start", "go run main.go", "Command to run to start the module")
	moduleClientCmd.PersistentFlags().StringVar(&PipelineUpRegexIgnore, "regex-ignore", "*.pb.go", "Regular expression for files to ignore")
	moduleClientCmd.PersistentFlags().StringVar(&PipelineDownModules, "modules-pull", "", "Comma-seperated list of the names of the modules to pull")
	moduleClientCmd.PersistentFlags().StringVar(&PipelineDownDirModules, "dir-pull", filepath.Join(os.TempDir(), "godibs", "pull", id), "Directory to pull the modules to")

	moduleCmd.AddCommand(moduleClientCmd)
}
