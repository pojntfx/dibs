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
			platformFromConfig := viper.GetString(PlatformKey)

			// Ignore if there are errors here, config and platforms might not be set (there is no hard dependency on the config)
			_ = ReadConfig(viper.GetString(DibsFileKey))
			platforms, _ := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)

			ignoreRegex := SyncClientIgnoreRegexPlaceholder
			if len(platforms) > 0 && viper.GetString(SyncClientPipelineUpRegexIgnoreKey) == SyncClientIgnoreRegexPlaceholder {
				ignoreRegex = platforms[0].Assets.CleanGlobs[0]
			}

			client := starters.Client{
				PipelineUpFileMod:      viper.GetString(SyncClientGoPipelineUpFileModKey),
				PipelineDownModules:    viper.GetString(SyncClientGoPipelineDownModulesKey),
				PipelineDownDirModules: viper.GetString(SyncClientGoPipelineDownDirModulesKey),
				PipelineUpBuildCommand: strings.Replace(viper.GetString(SyncClientPipelineUpBuildCommandKey), SyncClientPlatformPlaceholder, viper.GetString(PlatformKey), -1),
				PipelineUpStartCommand: strings.Replace(viper.GetString(SyncClientPipelineUpStartCommandKey), SyncClientPlatformPlaceholder, viper.GetString(PlatformKey), -1),
				PipelineUpTestCommand:  strings.Replace(viper.GetString(SyncClientPipelineUpTestCommandKey), SyncClientPlatformPlaceholder, viper.GetString(PlatformKey), -1),
				PipelineUpDirSrc:       viper.GetString(SyncClientPipelineUpDirSrcKey),
				PipelineUpDirPush:      viper.GetString(SyncClientPipelineUpDirPushKey),
				PipelineUpDirWatch:     viper.GetString(SyncClientPipelineUpDirWatchKey),
				PipelineUpRegexIgnore:  strings.Replace(viper.GetString(SyncClientPipelineUpRegexIgnoreKey), SyncClientIgnoreRegexPlaceholder, ignoreRegex, -1),

				RedisUrl:                  viper.GetString(SyncRedisUrlKey),
				RedisPrefix:               viper.GetString(SyncRedisPrefixKey),
				RedisPassword:             viper.GetString(SyncRedisPasswordKey),
				RedisSuffixUpRegistered:   SyncRedisSuffixUpRegistered,
				RedisSuffixUpUnRegistered: SyncRedisSuffixUpUnregistered,
				RedisSuffixUpTested:       SyncRedisSuffixUpTested,
				RedisSuffixUpBuilt:        SyncRedisSuffixUpBuilt,
				RedisSuffixUpStarted:      SyncRedisSuffixUpStarted,
				RedisSuffixUpPushed:       SyncRedisSuffixUpPushed,

				GitUpRemoteName:    SyncClientGitUpRemoteName,
				GitUpBaseURL:       viper.GetString(SyncClientGoGitBaseUrlKey),
				GitUpUserName:      SyncClientGitUpUserName,
				GitUpUserEmail:     SyncClientGitUpUserEmail,
				GitUpCommitMessage: SyncClientGitUpCommitMessageUpSynced,
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

		goGitBaseUrlFlag = strings.Replace(SyncClientGoGitBaseUrlKey, "_", "-", -1)

		pipelineUpDirSrcFlag   = strings.Replace(SyncClientPipelineUpDirSrcKey, "_", "-", -1)
		pipelineUpDirPushFlag  = strings.Replace(SyncClientPipelineUpDirPushKey, "_", "-", -1)
		pipelineUpDirWatchFlag = strings.Replace(SyncClientPipelineUpDirWatchKey, "_", "-", -1)

		goPipelineUpFileModFlag = strings.Replace(SyncClientGoPipelineUpFileModKey, "_", "-", -1)

		pipelineUpBuildCommandFlag = strings.Replace(SyncClientPipelineUpBuildCommandKey, "_", "-", -1)
		pipelineUpTestCommandFlag  = strings.Replace(SyncClientPipelineUpTestCommandKey, "_", "-", -1)
		pipelineUpStartCommandFlag = strings.Replace(SyncClientPipelineUpStartCommandKey, "_", "-", -1)

		pipelineUpRegexIgnoreFlag    = strings.Replace(SyncClientPipelineUpRegexIgnoreKey, "_", "-", -1)
		goPipelineDownModulesFlag    = strings.Replace(SyncClientGoPipelineDownModulesKey, "_", "-", -1)
		goPipelineDownDirModulesFlag = strings.Replace(SyncClientGoPipelineDownDirModulesKey, "_", "-", -1)

		id = uuid.New().String()
	)

	PipelineSyncClientCmd.PersistentFlags().StringVar(&goGitUpBaseUrl, goGitBaseUrlFlag, "http://localhost:32000/repos", `(--lang "`+LangGo+`" only) Base URL of the sync remote`)

	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpDirSrc, pipelineUpDirSrcFlag, ".", "Directory in which the source code of the module to push resides")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpDirPush, pipelineUpDirPushFlag, filepath.Join(os.TempDir(), "dibs", "push", id), "Temporary directory to put the module into before pushing")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpDirWatch, pipelineUpDirWatchFlag, ".", "Directory to watch for changes")

	PipelineSyncClientCmd.PersistentFlags().StringVar(&goPipelineUpFileMod, goPipelineUpFileModFlag, "go.mod", `(--lang "`+LangGo+`" only) Go module file of the module to push`)

	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpBuildCommand, pipelineUpBuildCommandFlag, os.Args[0]+" --platform "+SyncClientPlatformPlaceholder+" pipeline build assets", "Command to run to build the module. Infers the platform from the parent command by default")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpTestCommand, pipelineUpTestCommandFlag, os.Args[0]+" --platform "+SyncClientPlatformPlaceholder+" pipeline test unit lang", "Command to run to test the module. Infers the platform from the parent command by default")
	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpStartCommand, pipelineUpStartCommandFlag, os.Args[0]+" --platform "+SyncClientPlatformPlaceholder+" pipeline test integration assets", "Command to run to start the module. Infers the platform from the parent command by default")

	PipelineSyncClientCmd.PersistentFlags().StringVar(&pipelineUpRegexIgnore, pipelineUpRegexIgnoreFlag, SyncClientIgnoreRegexPlaceholder, "Regular expression for files to ignore. If a dibs configuration file exists, it will infer it from assets.cleanGlob")
	PipelineSyncClientCmd.PersistentFlags().StringVarP(&goPipelineDownModules, goPipelineDownModulesFlag, "g", "", `(--lang "`+LangGo+`" only) Comma-separated list of the names of the modules to pull`)
	PipelineSyncClientCmd.PersistentFlags().StringVar(&goPipelineDownDirModules, goPipelineDownDirModulesFlag, filepath.Join(os.TempDir(), "dibs", "pull", id), `(--lang "`+LangGo+`" only) Directory to pull the modules to`)

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(SyncClientGoGitBaseUrlKey, PipelineSyncClientCmd.PersistentFlags().Lookup(goGitBaseUrlFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(SyncClientPipelineUpDirSrcKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpDirSrcFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(SyncClientPipelineUpDirPushKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpDirPushFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(SyncClientPipelineUpDirWatchKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpDirWatchFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(SyncClientGoPipelineUpFileModKey, PipelineSyncClientCmd.PersistentFlags().Lookup(goPipelineUpFileModFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(SyncClientPipelineUpBuildCommandKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpBuildCommandFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(SyncClientPipelineUpTestCommandKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpTestCommandFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(SyncClientPipelineUpStartCommandKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpStartCommandFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(SyncClientPipelineUpRegexIgnoreKey, PipelineSyncClientCmd.PersistentFlags().Lookup(pipelineUpRegexIgnoreFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(SyncClientGoPipelineDownModulesKey, PipelineSyncClientCmd.PersistentFlags().Lookup(goPipelineDownModulesFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(SyncClientGoPipelineDownDirModulesKey, PipelineSyncClientCmd.PersistentFlags().Lookup(goPipelineDownDirModulesFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	viper.AutomaticEnv()

	PipelineSyncCmd.AddCommand(PipelineSyncClientCmd)
}
