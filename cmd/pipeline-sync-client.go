package cmd

import (
	"github.com/google/uuid"
	"github.com/pojntfx/dibs/pkg/starters"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var PipelineSyncClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start the module development client",
	Run: func(cmd *cobra.Command, args []string) {
		switch Lang {
		case LangGo:
			// Ignore if there are errors here, platforms might not be set (there is no hard dependency on the config)
			platforms, _ := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
			ignoreRegex := IgnoreRegexPlaceholder
			if len(platforms) > 0 {
				ignoreRegex = platforms[0].Assets.CleanGlob
			}

			client := starters.Client{
				PipelineUpFileMod:      GoPipelineUpFileMod,
				PipelineDownModules:    GoPipelineDownModules,
				PipelineDownDirModules: GoPipelineDownDirModules,
				PipelineUpBuildCommand: strings.Replace(PipelineUpBuildCommand, PlatformPlaceholder, Platform, -1),
				PipelineUpStartCommand: strings.Replace(PipelineUpStartCommand, PlatformPlaceholder, Platform, -1),
				PipelineUpTestCommand:  strings.Replace(PipelineUpTestCommand, PlatformPlaceholder, Platform, -1),
				PipelineUpDirSrc:       PipelineUpDirSrc,
				PipelineUpDirPush:      PipelineUpDirPush,
				PipelineUpDirWatch:     PipelineUpDirWatch,
				PipelineUpRegexIgnore:  strings.Replace(PipelineUpRegexIgnore, IgnoreRegexPlaceholder, ignoreRegex, -1),

				RedisUrl:                  RedisUrl,
				RedisPrefix:               RedisPrefix,
				RedisSuffixUpRegistered:   RedisSuffixUpRegistered,
				RedisSuffixUpUnRegistered: RedisSuffixUpUnregistered,
				RedisSuffixUpTested:       RedisSuffixUpTested,
				RedisSuffixUpBuilt:        RedisSuffixUpBuilt,
				RedisSuffixUpStarted:      RedisSuffixUpStarted,
				RedisSuffixUpPushed:       RedisSuffixUpPushed,

				GitUpRemoteName:    GitUpRemoteName,
				GitUpBaseURL:       viper.GetString(GoGitBaseUrlKey),
				GitUpUserName:      GitUpUserName,
				GitUpUserEmail:     GitUpUserEmail,
				GitUpCommitMessage: GitUpCommitMessage,
			}

			client.Start()
		}
	},
}

var (
	GoGitUpBaseUrl string

	PipelineUpDirSrc       string
	PipelineUpDirPush      string
	PipelineUpDirWatch     string
	GoPipelineUpFileMod    string
	PipelineUpBuildCommand string
	PipelineUpTestCommand  string
	PipelineUpStartCommand string
	PipelineUpRegexIgnore  string

	GoPipelineDownModules    string
	GoPipelineDownDirModules string
)

const (
	GitUpCommitMessage = "up_synced"
	GitUpRemoteName    = "dibs-sync"
	GitUpUserName      = "dibs-syncer"
	GitUpUserEmail     = "dibs-syncer@pojtinger.space"

	PlatformPlaceholder    = "[infer]"
	IgnoreRegexPlaceholder = "[infer]"

	EnvPrefix = "dibs"

	GoGitBaseUrlFlag = LangGo + "-git-base-url"
	GoGitBaseUrlKey  = LangGo + "_git_base_url"
)

func init() {
	id := uuid.New().String()

	PipelineSyncClientCmd.PersistentFlags().StringVar(&GoGitUpBaseUrl, GoGitBaseUrlFlag, "http://localhost:35000/repos", `(--lang "`+LangGo+`" only) Base URL of the sync remote`)

	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpDirSrc, "dir-src", ".", "Directory in which the source code of the module to push resides")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpDirPush, "dir-push", filepath.Join(os.TempDir(), "dibs", "push", id), "Temporary directory to put the module into before pushing")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpDirWatch, "dir-watch", ".", "Directory to watch for changes")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&GoPipelineUpFileMod, LangGo+"-modules-file", "go.mod", `(--lang "`+LangGo+`" only) Go module file of the module to push`)
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpBuildCommand, "cmd-build", os.Args[0]+" --platform "+PlatformPlaceholder+" pipeline build assets", "Command to run to build the module. Infers the platform from the parent command by default")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpTestCommand, "cmd-test", os.Args[0]+" --platform "+PlatformPlaceholder+" pipeline test unit lang", "Command to run to test the module. Infers the platform from the parent command by default")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpStartCommand, "cmd-start", os.Args[0]+" --platform "+PlatformPlaceholder+" pipeline test integration assets", "Command to run to start the module. Infers the platform from the parent command by default")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&PipelineUpRegexIgnore, "regex-ignore", IgnoreRegexPlaceholder, "Regular expression for files to ignore. If a dibs configuration file exists, it will infer it from assets.cleanGlob")
	PipelineSyncClientCmd.PersistentFlags().StringVarP(&GoPipelineDownModules, LangGo+"-modules-pull", "g", "", `(--lang "`+LangGo+`" only) Comma-separated list of the names of the modules to pull`)
	PipelineSyncClientCmd.PersistentFlags().StringVar(&GoPipelineDownDirModules, LangGo+"-dir-pull", filepath.Join(os.TempDir(), "dibs", "pull", id), `(--lang "`+LangGo+`" only) Directory to pull the modules to`)

	viper.SetEnvPrefix(EnvPrefix)

	viper.BindPFlag(GoGitBaseUrlKey, PipelineSyncClientCmd.PersistentFlags().Lookup(GoGitBaseUrlFlag))

	viper.AutomaticEnv()

	PipelineSyncCmd.AddCommand(PipelineSyncClientCmd)
}
