package cmd

import (
	"github.com/pojntfx/godibs/pkg/config"
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/pojntfx/godibs/pkg/workers"
	"github.com/spf13/cobra"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"path/filepath"
	"strconv"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		// Connect to Redis
		redis := utils.Redis{
			Addr:   REDIS_URL,
			Prefix: REDIS_PREFIX,
		}
		redis.Connect()

		// Build the configuration
		httpPort, err := strconv.ParseInt(config.GIT_SERVER_HTTP_PORT, 0, 64)
		if err != nil {
			log.Fatal("Error", rz.String("System", "Server"), rz.Err(err))
		}
		reposDirWithHTTPPathPrefix := filepath.Join(config.GIT_SERVER_REPOS_DIR, config.GIT_SERVER_HTTP_PATH)

		// Setup workers
		httpWorker := &workers.GitHTTPWorker{
			ReposDir:       config.GIT_SERVER_REPOS_DIR,
			HTTPPathPrefix: config.GIT_SERVER_HTTP_PATH,
			Port:           int(httpPort),
		}

		repoWorkerUpdate, repoWorkerDeleteOnly := &workers.GitRepoWorker{
			ReposDir:    reposDirWithHTTPPathPrefix,
			DeleteOnly:  false,
			Redis:       redis,
			RedisSuffix: REDIS_SUFFIX_UP_REGISTERED,
		}, &workers.GitRepoWorker{
			ReposDir:    reposDirWithHTTPPathPrefix,
			DeleteOnly:  true,
			Redis:       redis,
			RedisSuffix: REDIS_SUFFIX_UP_UNREGISTERED,
		}

		// Create error channels
		httpWorkerErrors, repoWorkerUpdateErrors, repoWorkerDeleteOnlyErrors := make(chan error, 0), make(chan error, 0), make(chan error, 0)

		// Create event channels
		httpWorkerEvents, repoWorkerUpdateEvents, repoWorkerDeleteOnlyEvents := make(chan utils.Event, 0), make(chan utils.Event, 0), make(chan utils.Event, 0)

		// Start workers
		go httpWorker.Start(httpWorkerErrors, httpWorkerEvents)
		go repoWorkerUpdate.Start(repoWorkerUpdateErrors, repoWorkerUpdateEvents)
		go repoWorkerDeleteOnly.Start(repoWorkerDeleteOnlyErrors, repoWorkerDeleteOnlyEvents)

		// Start main loop
		for {
			select {
			// If there are errors, log the erros and exit
			case err := <-httpWorkerErrors:
				log.Fatal("Error", rz.String("System", "GitHTTPWorker"), rz.Err(err))
			case err := <-repoWorkerUpdateErrors:
				log.Fatal("Error", rz.String("System", "GitRepoWorker"), rz.Bool("DeleteOnly", repoWorkerUpdate.DeleteOnly), rz.Err(err))
			case err := <-repoWorkerDeleteOnlyErrors:
				log.Fatal("Error", rz.String("System", "GitRepoWorker"), rz.Bool("DeleteOnly", repoWorkerDeleteOnly.DeleteOnly), rz.Err(err))

			// If there are events, log them
			case event := <-httpWorkerEvents:
				switch event.Code {
				case 0:
					log.Info("Started", rz.String("System", "GitHTTPWorker"), rz.String("EventMessage", event.Message), rz.String("ReposDir", httpWorker.ReposDir), rz.String("HTTPPathPrefix", httpWorker.HTTPPathPrefix), rz.Int("Port", httpWorker.Port))
				case 1:
					log.Info("Request", rz.String("System", "GitHTTPWorker"), rz.String("EventMessage", event.Message))
				case 2:
					log.Info("Stopped", rz.String("System", "GitHTTPWorker"), rz.String("EventMessage", event.Message))
					return
				default:
					log.Fatal("Unknown event code", rz.String("System", "GitHTTPWorker"), rz.Int("EventCode", event.Code), rz.String("StatusMessage", event.Message))
				}
			case event := <-repoWorkerUpdateEvents:
				switch event.Code {
				case 0:
					log.Info("Started", rz.String("System", "GitRepoWorker"), rz.String("EventMessage", event.Message), rz.String("ReposDir", repoWorkerUpdate.ReposDir), rz.Bool("DeleteOnly", repoWorkerUpdate.DeleteOnly), rz.String("RedisSuffix", repoWorkerUpdate.RedisSuffix))
				case 1:
					log.Info("Update", rz.String("System", "GitRepoWorker"), rz.Bool("DeleteOnly", repoWorkerUpdate.DeleteOnly), rz.String("EventMessage", event.Message))
				case 2:
					log.Info("Stopped", rz.String("System", "GitRepoWorker"), rz.Bool("DeleteOnly", repoWorkerUpdate.DeleteOnly), rz.String("EventMessage", event.Message))
					return
				default:
					log.Fatal("Unknown event code", rz.String("System", "GitRepoWorker"), rz.Bool("DeleteOnly", repoWorkerUpdate.DeleteOnly), rz.Int("EventCode", event.Code), rz.String("StatusMessage", event.Message))
				}
			case event := <-repoWorkerDeleteOnlyEvents:
				switch event.Code {
				case 0:
					log.Info("Started", rz.String("System", "GitRepoWorker"), rz.String("EventMessage", event.Message), rz.String("ReposDir", repoWorkerDeleteOnly.ReposDir), rz.Bool("DeleteOnly", repoWorkerDeleteOnly.DeleteOnly), rz.String("RedisSuffix", repoWorkerDeleteOnly.RedisSuffix))
				case 1:
					log.Info("Deletion", rz.String("System", "GitRepoWorker"), rz.Bool("DeleteOnly", repoWorkerDeleteOnly.DeleteOnly), rz.String("EventMessage", event.Message))
				case 2:
					log.Info("Stopped", rz.String("System", "GitRepoWorker"), rz.Bool("DeleteOnly", repoWorkerDeleteOnly.DeleteOnly), rz.String("EventMessage", event.Message))
					return
				default:
					log.Fatal("Unknown event code", rz.String("System", "GitRepoWorker"), rz.Bool("DeleteOnly", repoWorkerUpdate.DeleteOnly), rz.Int("EventCode", event.Code), rz.String("StatusMessage", event.Message))
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
