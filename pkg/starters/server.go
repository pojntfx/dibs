package starters

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/pojntfx/dibs/pkg/workers"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"path/filepath"
	"strconv"
)

// Server is a server for the sync client
type Server struct {
	ServerReposDir string // Directory in which the Git repos should be stored
	ServerHTTPPort string // Port on which the Git repos should be served
	ServerHTTPPath string // HTTP path prefix for the served Git repos

	RedisUrl                  string // URL of the Redis instance to use
	RedisPrefix               string // Redis channel prefix
	RedisPassword             string // Redis password
	RedisSuffixUpRegistered   string // Redis channel suffix for "module registered" messages
	RedisSuffixUpUnRegistered string // Redis channel suffix for "module unregistered" messages
}

// Start starts the sync server
func (server *Server) Start() {
	// Connect to Redis
	redis := utils.Redis{
		Addr:     server.RedisUrl,
		Prefix:   server.RedisPrefix,
		Password: server.RedisPassword,
	}
	redis.Connect()

	// Build the configuration
	httpPort, err := strconv.ParseInt(server.ServerHTTPPort, 0, 64)
	if err != nil {
		utils.LogErrorFatal("Error", err)
	}
	reposDirWithHTTPPathPrefix := filepath.Join(server.ServerReposDir, server.ServerHTTPPath)

	// Setup workers
	httpWorker := &workers.GitHTTPWorker{
		ReposDir:       server.ServerReposDir,
		HTTPPathPrefix: server.ServerHTTPPath,
		Port:           int(httpPort),
	}

	repoWorkerUpdate, repoWorkerDeleteOnly := &workers.GitRepoWorker{
		ReposDir:    reposDirWithHTTPPathPrefix,
		DeleteOnly:  false,
		Redis:       redis,
		RedisSuffix: server.RedisSuffixUpRegistered,
	}, &workers.GitRepoWorker{
		ReposDir:    reposDirWithHTTPPathPrefix,
		DeleteOnly:  true,
		Redis:       redis,
		RedisSuffix: server.RedisSuffixUpUnRegistered,
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
		// If there are errors, log the errors and exit
		case err := <-httpWorkerErrors:
			log.Fatal("Error", rz.String("system", "GitHTTPWorker"), rz.Err(err))
		case err := <-repoWorkerUpdateErrors:
			log.Fatal("Error", rz.String("system", "GitRepoWorker"), rz.Bool("deleteOnly", repoWorkerUpdate.DeleteOnly), rz.Err(err))
		case err := <-repoWorkerDeleteOnlyErrors:
			log.Fatal("Error", rz.String("system", "GitRepoWorker"), rz.Bool("deleteOnly", repoWorkerDeleteOnly.DeleteOnly), rz.Err(err))

		// If there are events, log them
		case event := <-httpWorkerEvents:
			switch event.Code {
			case 0:
				log.Info("Started", rz.String("system", "GitHTTPWorker"), rz.String("eventMessage", event.Message), rz.String("repositoriesDir", httpWorker.ReposDir), rz.String("HTTPPathPrefix", httpWorker.HTTPPathPrefix), rz.Int("Port", httpWorker.Port))
			case 1:
				log.Info("Request", rz.String("system", "GitHTTPWorker"), rz.String("eventMessage", event.Message))
			case 2:
				log.Info("Stopped", rz.String("system", "GitHTTPWorker"), rz.String("eventMessage", event.Message))
				return
			default:
				log.Fatal("Unknown event code", rz.String("system", "GitHTTPWorker"), rz.Int("eventCode", event.Code), rz.String("statusMessage", event.Message))
			}
		case event := <-repoWorkerUpdateEvents:
			switch event.Code {
			case 0:
				log.Info("Started", rz.String("system", "GitRepoWorker"), rz.String("eventMessage", event.Message), rz.String("repositoriesDir", repoWorkerUpdate.ReposDir), rz.Bool("deleteOnly", repoWorkerUpdate.DeleteOnly), rz.String("redisSuffix", repoWorkerUpdate.RedisSuffix))
			case 1:
				log.Info("Update", rz.String("system", "GitRepoWorker"), rz.Bool("deleteOnly", repoWorkerUpdate.DeleteOnly), rz.String("eventMessage", event.Message))
			case 2:
				log.Info("Stopped", rz.String("system", "GitRepoWorker"), rz.Bool("deleteOnly", repoWorkerUpdate.DeleteOnly), rz.String("eventMessage", event.Message))
				return
			default:
				log.Fatal("Unknown event code", rz.String("system", "GitRepoWorker"), rz.Bool("deleteOnly", repoWorkerUpdate.DeleteOnly), rz.Int("eventCode", event.Code), rz.String("statusMessage", event.Message))
			}
		case event := <-repoWorkerDeleteOnlyEvents:
			switch event.Code {
			case 0:
				log.Info("Started", rz.String("system", "GitRepoWorker"), rz.String("eventMessage", event.Message), rz.String("repositoriesDir", repoWorkerDeleteOnly.ReposDir), rz.Bool("deleteOnly", repoWorkerDeleteOnly.DeleteOnly), rz.String("redisSuffix", repoWorkerDeleteOnly.RedisSuffix))
			case 1:
				log.Info("Deletion", rz.String("system", "GitRepoWorker"), rz.Bool("deleteOnly", repoWorkerDeleteOnly.DeleteOnly), rz.String("eventMessage", event.Message))
			case 2:
				log.Info("Stopped", rz.String("system", "GitRepoWorker"), rz.Bool("deleteOnly", repoWorkerDeleteOnly.DeleteOnly), rz.String("eventMessage", event.Message))
				return
			default:
				log.Fatal("Unknown event code", rz.String("system", "GitRepoWorker"), rz.Bool("deleteOnly", repoWorkerUpdate.DeleteOnly), rz.Int("eventCode", event.Code), rz.String("statusMessage", event.Message))
			}
		}
	}
}
