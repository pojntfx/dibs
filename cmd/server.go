package cmd

import (
	"github.com/pojntfx/godibs/pkg/starters"
	"github.com/spf13/cobra"
)

var (
	GIT_SERVER_REPOS_DIR string
	GIT_SERVER_HTTP_PORT string
	GIT_SERVER_HTTP_PATH string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		server := starters.Server{
			ServerReposDir:            GIT_SERVER_REPOS_DIR,
			ServerHTTPPort:            GIT_SERVER_HTTP_PORT,
			ServerHTTPPath:            GIT_SERVER_HTTP_PATH,
			RedisUrl:                  REDIS_URL,
			RedisPrefix:               REDIS_PREFIX,
			RedisSuffixUpRegistered:   REDIS_SUFFIX_UP_REGISTERED,
			RedisSuffixUpUnRegistered: REDIS_SUFFIX_UP_UNREGISTERED,
		}

		server.Start()
	},
}

func init() {
	serverCmd.PersistentFlags().StringVar(&GIT_SERVER_REPOS_DIR, "dir-repos", "/tmp/serverrepos", "Directory in which the Git repos should be stored")
	serverCmd.PersistentFlags().StringVar(&GIT_SERVER_HTTP_PORT, "port", "25000", "Port on which the Git repos should be served")
	serverCmd.PersistentFlags().StringVar(&GIT_SERVER_HTTP_PATH, "path", "/repos", "HTTP path on which the Git repos should be served")

	rootCmd.AddCommand(serverCmd)
}
