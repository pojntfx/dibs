package cmd

import (
	"github.com/google/uuid"
	"github.com/pojntfx/dibs/pkg/starters"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"os"
	"path/filepath"
	"strings"
)

var PipelineSyncClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start the module development client",
	Run: func(cmd *cobra.Command, args []string) {
		switch viper.GetString(LangKey) {
		case LangGo:
			// Ignore if there are errors here, platforms might not be set (there is no hard dependency on the config)
			platformFromConfig := viper.GetString(PlatformKey)

			platforms, _ := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
			ignoreRegex := IgnoreRegexPlaceholder
			if len(platforms) > 0 {
				ignoreRegex = platforms[0].Assets.CleanGlob
			}

			client := starters.Client{
				PipelineUpFileMod:      viper.GetString(GoPipelineUpFileModKey),
				PipelineDownModules:    viper.GetString(GoPipelineDownModulesKey),
				PipelineDownDirModules: viper.GetString(GoPipelineDownDirModulesKey),
				PipelineUpBuildCommand: strings.Replace(viper.GetString(PipelineUpBuildCommandKey), PlatformPlaceholder, viper.GetString(PlatformKey), -1),
				PipelineUpStartCommand: strings.Replace(viper.GetString(PipelineUpStartCommandKey), PlatformPlaceholder, viper.GetString(PlatformKey), -1),
				PipelineUpTestCommand:  strings.Replace(viper.GetString(PipelineUpTestCommandKey), PlatformPlaceholder, viper.GetString(PlatformKey), -1),
				PipelineUpDirSrc:       viper.GetString(PipelineUpDirSrcKey),
				PipelineUpDirPush:      viper.GetString(PipelineUpDirPushKey),
				PipelineUpDirWatch:     viper.GetString(PipelineUpDirWatchKey),
				PipelineUpRegexIgnore:  strings.Replace(viper.GetString(PipelineUpRegexIgnoreKey), IgnoreRegexPlaceholder, ignoreRegex, -1),

				RedisUrl:                  viper.GetString(RedisUrlKey),
				RedisPrefix:               viper.GetString(RedisPrefixKey),
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

func init() {
	var (
		goGitUpBaseUrl string

		pipelineUpDirSrc   string
		pipelineUpDirPush  string
		pipelineUpDirWatch string

		goPipelineUpFileMod string

		pipelineUpBuildCommand string
		pipelineUpTestCommand  string
		pipelineUpStartCommand string

		pipelineUpRegexIgnore    string
		goPipelineDownModules    string
		goPipelineDownDirModules string

		goGitBaseUrlFlag = strings.Replace(GoGitBaseUrlKey, "_", "-", -1)

		pipelineUpDirSrcFlag   = strings.Replace(PipelineUpDirSrcKey, "_", "-", -1)
		pipelineUpDirPushFlag  = strings.Replace(PipelineUpDirPushKey, "_", "-", -1)
		pipelineUpDirWatchFlag = strings.Replace(PipelineUpDirWatchKey, "_", "-", -1)

		goPipelineUpFileModFlag = strings.Replace(GoPipelineUpFileModKey, "_", "-", -1)

		pipelineUpBuildCommandFlag = strings.Replace(PipelineUpBuildCommandKey, "_", "-", -1)
		pipelineUpTestCommandFlag  = strings.Replace(PipelineUpTestCommandKey, "_", "-", -1)
		pipelineUpStartCommandFlag = strings.Replace(PipelineUpStartCommandKey, "_", "-", -1)

		pipelineUpRegexIgnoreFlag    = strings.Replace(PipelineUpRegexIgnoreKey, "_", "-", -1)
		goPipelineDownModulesFlag    = strings.Replace(GoPipelineDownModulesKey, "_", "-", -1)
		goPipelineDownDirModulesFlag = strings.Replace(GoPipelineDownDirModulesKey, "_", "-", -1)

		id = uuid.New().String()
	)

	PipelineSyncClientCmd.PersistentFlags().StringVar(&goGitUpBaseUrl, goGitBaseUrlFlag, "http://localhost:35000/repos", `(--lang "`+LangGo+`" only) Base URL of the sync remote`)

	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpDirSrc, pipelineUpDirSrcFlag, ".", "Directory in which the source code of the module to push resides")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpDirPush, pipelineUpDirPushFlag, filepath.Join(os.TempDir(), "dibs", "push", id), "Temporary directory to put the module into before pushing")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpDirWatch, pipelineUpDirWatchFlag, ".", "Directory to watch for changes")

	PipelineSyncClientCmd.PersistentFlags().StringVar(&goPipelineUpFileMod, goPipelineUpFileModFlag, "go.mod", `(--lang "`+LangGo+`" only) Go module file of the module to push`)

	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpBuildCommand, pipelineUpBuildCommandFlag, os.Args[0]+" --platform "+PlatformPlaceholder+" pipeline build assets", "Command to run to build the module. Infers the platform from the parent command by default")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpTestCommand, pipelineUpTestCommandFlag, os.Args[0]+" --platform "+PlatformPlaceholder+" pipeline test unit lang", "Command to run to test the module. Infers the platform from the parent command by default")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpStartCommand, pipelineUpStartCommandFlag, os.Args[0]+" --platform "+PlatformPlaceholder+" pipeline test integration assets", "Command to run to start the module. Infers the platform from the parent command by default")

	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpRegexIgnore, pipelineUpRegexIgnoreFlag, IgnoreRegexPlaceholder, "Regular expression for files to ignore. If a dibs configuration file exists, it will infer it from assets.cleanGlob")
	PipelineSyncClientCmd.PersistentFlags().StringVarP(&goPipelineDownModules, goPipelineDownModulesFlag, "g", "", `(--lang "`+LangGo+`" only) Comma-separated list of the names of the modules to pull`)
	PipelineSyncClientCmd.PersistentFlags().StringVar(&goPipelineDownDirModules, goPipelineDownDirModulesFlag, filepath.Join(os.TempDir(), "dibs", "pull", id), `(--lang "`+LangGo+`" only) Directory to pull the modules to`)

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(GoGitBaseUrlKey, PipelineSyncClientCmd.PersistentFlags().Lookup(goGitBaseUrlFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(PipelineUpDirSrcKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpDirSrcFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(PipelineUpDirPushKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpDirPushFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(PipelineUpDirWatchKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpDirWatchFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(GoPipelineUpFileModKey, PipelineSyncClientCmd.PersistentFlags().Lookup(goPipelineUpFileModFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(PipelineUpBuildCommandKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpBuildCommandFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(PipelineUpTestCommandKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpTestCommandFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(PipelineUpStartCommandKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpStartCommandFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(PipelineUpRegexIgnoreKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpRegexIgnoreFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(GoPipelineDownModulesKey, PipelineSyncClientCmd.PersistentFlags().Lookup(goPipelineDownModulesFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(GoPipelineDownDirModulesKey, PipelineSyncClientCmd.PersistentFlags().Lookup(goPipelineDownDirModulesFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	viper.AutomaticEnv()

	PipelineSyncCmd.AddCommand(PipelineSyncClientCmd)
}
