package cmd

import (
	"github.com/google/uuid"
	"github.com/pojntfx/godibs/pkg/starters"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	GIT_SERVER_REPOS_DIR string
	GIT_SERVER_HTTP_PORT string
	GIT_SERVER_HTTP_PATH string
)

// moduleServerCmd ist the command to start the server
var moduleServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		server := starters.Server{
			ServerReposDir: GIT_SERVER_REPOS_DIR,
			ServerHTTPPort: GIT_SERVER_HTTP_PORT,
			ServerHTTPPath: GIT_SERVER_HTTP_PATH,

			RedisUrl:                  REDIS_URL,
			RedisPrefix:               REDIS_PREFIX,
			RedisSuffixUpRegistered:   REDIS_SUFFIX_UP_REGISTERED,
			RedisSuffixUpUnRegistered: REDIS_SUFFIX_UP_UNREGISTERED,
		}

		server.Start()
	},
}

// init maps the flags to the config
func init() {
	id := uuid.New().String()

	moduleServerCmd.PersistentFlags().StringVar(&GIT_SERVER_REPOS_DIR, "dir-repos", filepath.Join(os.TempDir(), "godibs", "gitrepos", id), "Directory in which the Git repos should be stored")
	moduleServerCmd.PersistentFlags().StringVar(&GIT_SERVER_HTTP_PORT, "port", "25000", "Port on which the Git repos should be served")
	moduleServerCmd.PersistentFlags().StringVar(&GIT_SERVER_HTTP_PATH, "path", "/repos", "HTTP path prefix for the served Git repos")

	moduleCmd.AddCommand(moduleServerCmd)
}
