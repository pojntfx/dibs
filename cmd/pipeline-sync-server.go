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

var PipelineSyncServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the module development server",
	Run: func(cmd *cobra.Command, args []string) {
		switch viper.GetString(LangKey) {
		case LangGo:
			server := starters.Server{
				ServerReposDir: viper.GetString(GitServerReposDirKey),
				ServerHTTPPort: viper.GetString(GitServerHttpPortKey),
				ServerHTTPPath: viper.GetString(GitServerHttpPathKey),

				RedisUrl:                  viper.GetString(RedisUrlKey),
				RedisPrefix:               viper.GetString(RedisPrefixKey),
				RedisSuffixUpRegistered:   RedisSuffixUpRegistered,
				RedisSuffixUpUnRegistered: RedisSuffixUpUnregistered,
			}

			server.Start()
		}
	},
}

const (
	RedisSuffixUpBuilt        = "up_built"
	RedisSuffixUpTested       = "up_tested"
	RedisSuffixUpStarted      = "up_started"
	RedisSuffixUpRegistered   = "up_registered"
	RedisSuffixUpUnregistered = "up_unregistered"
	RedisSuffixUpPushed       = "up_pushed"

	GitServerReposDirKey = LangGo + "_dir_repos"
	GitServerHttpPortKey = LangGo + "-port"
	GitServerHttpPathKey = LangGo + "-path"
)

func init() {
	var (
		gitServerReposDir string
		gitServerHttpPort string
		gitServerHttpPath string

		gitServerReposDirFlag = strings.Replace(GitServerReposDirKey, "_", "-", -1)
		gitServerHttpPortFlag = strings.Replace(GitServerHttpPortKey, "_", "-", -1)
		gitServerHttpPathFlag = strings.Replace(GitServerHttpPathKey, "_", "-", -1)

		id = uuid.New().String()
	)

	PipelineSyncServerCmd.PersistentFlags().StringVar(&gitServerReposDir, gitServerReposDirFlag, filepath.Join(os.TempDir(), "dibs", "gitrepos", id), `(--lang "`+LangGo+`" only) Directory in which the Git repos should be stored`)
	PipelineSyncServerCmd.PersistentFlags().StringVar(&gitServerHttpPort, gitServerHttpPortFlag, "35000", `(--lang "`+LangGo+`" only) Port on which the Git repos should be served`)
	PipelineSyncServerCmd.PersistentFlags().StringVar(&gitServerHttpPath, gitServerHttpPathFlag, "/repos", `(--lang "`+LangGo+`" only) HTTP path prefix for the served Git repos`)

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(GitServerReposDirKey, PipelineSyncServerCmd.PersistentFlags().Lookup(gitServerReposDirFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(GitServerHttpPortKey, PipelineSyncServerCmd.PersistentFlags().Lookup(gitServerHttpPortFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(GitServerHttpPathKey, PipelineSyncServerCmd.PersistentFlags().Lookup(gitServerHttpPathFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	viper.AutomaticEnv()

	PipelineSyncCmd.AddCommand(PipelineSyncServerCmd)
}
