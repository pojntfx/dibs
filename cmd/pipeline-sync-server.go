package cmd

import (
	"github.com/google/uuid"
	"github.com/pojntfx/dibs/pkg/starters"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var PipelineSyncServerCmd = &cobra.Command{
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

var (
	GitServerReposDir string
	GitServerHttpPort string
	GitServerHttpPath string
)

const (
	RedisSuffixUpBuilt        = "up_built"
	RedisSuffixUpTested       = "up_tested"
	RedisSuffixUpStarted      = "up_started"
	RedisSuffixUpRegistered   = "up_registered"
	RedisSuffixUpUnregistered = "up_unregistered"
	RedisSuffixUpPushed       = "up_pushed"
)

func init() {
	id := uuid.New().String()

	PipelineSyncServerCmd.PersistentFlags().StringVar(&GitServerReposDir, "dir-repos", filepath.Join(os.TempDir(), "dibs", "gitrepos", id), "Directory in which the Git repos should be stored")
	PipelineSyncServerCmd.PersistentFlags().StringVar(&GitServerHttpPort, "port", "25000", "Port on which the Git repos should be served")
	PipelineSyncServerCmd.PersistentFlags().StringVar(&GitServerHttpPath, "path", "/repos", "HTTP path prefix for the served Git repos")

	PipelineSyncCmd.AddCommand(PipelineSyncServerCmd)
}
