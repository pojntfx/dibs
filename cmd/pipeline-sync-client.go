package cmd

import (
	"github.com/google/uuid"
	"github.com/pojntfx/dibs/pkg/starters"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var PipelineSyncClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start the module development client",
	Run: func(cmd *cobra.Command, args []string) {
		switch Lang {
		case LangGo:
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
		}
	},
}

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
	GitUpRemoteName    = "dibs-sync"
	GitUpUserName      = "dibs-syncer"
	GitUpUserEmail     = "dibs-syncer@pojtinger.space"
)

func init() {
	id := uuid.New().String()

	PipelineSyncClientCmd.PersistentFlags().StringVar(&GitUpBaseUrl, LangGo+"-git-base-url", "http://localhost:35000/repos", `(--lang "`+LangGo+`" only) Base URL of the sync remote`)

	platformToUse := PlatformDefault
	if Platform != PlatformDefault && Platform != "" {
		platformToUse = Platform
	} else if Platform == "" {
		platformToUse = "[inherited]"
	}

	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpDirSrc, "dir-src", ".", "Directory in which the source code of the module to push resides")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpDirPush, "dir-push", filepath.Join(os.TempDir(), "dibs", "push", id), "Temporary directory to put the module into before pushing")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpDirWatch, "dir-watch", ".", "Directory to watch for changes")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpFileMod, LangGo+"-modules-file", "go.mod", `(--lang "`+LangGo+`" only) Go module file of the module to push`)
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpBuildCommand, "cmd-build", os.Args[0]+" --platform "+platformToUse+" pipeline build assets", "Command to run to build the module")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpTestCommand, "cmd-test", os.Args[0]+" --platform "+platformToUse+" pipeline test unit lang", "Command to run to test the module")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpStartCommand, "cmd-start", os.Args[0]+" --platform "+platformToUse+" pipeline test integration assets", "Command to run to start the module")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpRegexIgnore, LangGo+"-regex-ignore", "*.pb.go", `(--lang "`+LangGo+`" only) Regular expression for files to ignore`)
	PipelineSyncClientCmd.PersistentFlags().StringVarP(&PipelineDownModules, LangGo+"-modules-pull", "g", "", `(--lang "`+LangGo+`" only) Comma-separated list of the names of the modules to pull`)
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineDownDirModules, LangGo+"-dir-pull", filepath.Join(os.TempDir(), "dibs", "pull", id), `(--lang "`+LangGo+`" only) Directory to pull the modules to`)

	PipelineSyncCmd.AddCommand(PipelineSyncClientCmd)
}
