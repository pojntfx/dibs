package cmd

import (
	"github.com/google/uuid"
	"github.com/pojntfx/dibs/pkg/starters"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var PipelineSyncServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the module development server",
	Run: func(cmd *cobra.Command, args []string) {
		switch viper.GetString(LangKey) {
		case LangGo:
			server := starters.Server{
				ServerReposDir: viper.GetString(SyncServerGitServerReposDirKey),
				ServerHTTPPort: viper.GetString(SyncServerGitServerHttpPortKey),
				ServerHTTPPath: viper.GetString(SyncServerGitServerHttpPathKey),

				RedisUrl:                  viper.GetString(SyncRedisUrlKey),
				RedisPrefix:               viper.GetString(SyncRedisPrefixKey),
				RedisPassword:             viper.GetString(SyncRedisPasswordKey),
				RedisSuffixUpRegistered:   SyncRedisSuffixUpRegistered,
				RedisSuffixUpUnRegistered: SyncRedisSuffixUpUnregistered,
			}

			server.Start()
		}
	},
}

func init() {
	var (
		gitServerReposDir string
		gitServerHttpPort string
		gitServerHttpPath string

		gitServerReposDirFlag = strings.Replace(strings.Replace(SyncServerGitServerReposDirKey, SyncKeyPrefix, "", -1), "_", "-", -1)
		gitServerHttpPortFlag = strings.Replace(strings.Replace(SyncServerGitServerHttpPortKey, SyncKeyPrefix, "", -1), "_", "-", -1)
		gitServerHttpPathFlag = strings.Replace(strings.Replace(SyncServerGitServerHttpPathKey, SyncKeyPrefix, "", -1), "_", "-", -1)

		id = uuid.New().String()
	)

	PipelineSyncServerCmd.PersistentFlags().StringVar(&gitServerReposDir, gitServerReposDirFlag, filepath.Join(os.TempDir(), "dibs", "gitRepos", id), `(--lang "`+LangGo+`" only) Directory in which the Git repos should be stored`)
	PipelineSyncServerCmd.PersistentFlags().StringVar(&gitServerHttpPort, gitServerHttpPortFlag, "32000", `(--lang "`+LangGo+`" only) Port on which the Git repos should be served`)
	PipelineSyncServerCmd.PersistentFlags().StringVar(&gitServerHttpPath, gitServerHttpPathFlag, "/repos", `(--lang "`+LangGo+`" only) HTTP path prefix for the served Git repos`)

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(SyncServerGitServerReposDirKey, PipelineSyncServerCmd.PersistentFlags().Lookup(gitServerReposDirFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(SyncServerGitServerHttpPortKey, PipelineSyncServerCmd.PersistentFlags().Lookup(gitServerHttpPortFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(SyncServerGitServerHttpPathKey, PipelineSyncServerCmd.PersistentFlags().Lookup(gitServerHttpPathFlag)); err != nil {
		utils.CmdLogErrorCouldNotBindFlag(err)
	}

	viper.AutomaticEnv()

	PipelineSyncCmd.AddCommand(PipelineSyncServerCmd)
}
