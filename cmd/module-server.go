package cmd

import (
	"github.com/google/uuid"
	"github.com/pojntfx/godibs/pkg/starters"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	GitServerReposDir string
	GitServerHttpPort string
	GitServerHttpPath string
)

// moduleServerCmd ist the command to start the server
var moduleServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the module development server",
	Run: func(cmd *cobra.Command, args []string) {
		server := starters.Server{
			ServerReposDir: GitServerReposDir,
			ServerHTTPPort: GitServerHttpPort,
			ServerHTTPPath: GitServerHttpPath,

			RedisUrl:                  RedisUrl,
			RedisPrefix:               RedisPrefix,
			RedisSuffixUpRegistered:   RedisSuffixUpRegistered,
			RedisSuffixUpUnRegistered: RedisSuffixUpUnregistered,
		}

		server.Start()
	},
}

// init maps the flags to the config
func init() {
	id := uuid.New().String()

	moduleServerCmd.PersistentFlags().StringVar(&GitServerReposDir, "dir-repos", filepath.Join(os.TempDir(), "godibs", "gitrepos", id), "Directory in which the Git repos should be stored")
	moduleServerCmd.PersistentFlags().StringVar(&GitServerHttpPort, "port", "25000", "Port on which the Git repos should be served")
	moduleServerCmd.PersistentFlags().StringVar(&GitServerHttpPath, "path", "/repos", "HTTP path prefix for the served Git repos")

	moduleCmd.AddCommand(moduleServerCmd)
}
