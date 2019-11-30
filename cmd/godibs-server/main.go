package main

import (
	redis "github.com/go-redis/redis/v7"
	"github.com/pojntfx/godibs/pkg/config"
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/pojntfx/godibs/pkg/workers"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	git "gopkg.in/src-d/go-git.v4"
	"os"
	"strconv"
	"sync"
)

func main() {
	redisClient := utils.NewRedisClient(config.REDIS_URL)

	httpPort, err := strconv.ParseInt(config.GIT_HTTP_PORT, 0, 64)
	if err != nil {
		panic(err)
	}

	httpWorker := &workers.GitHTTPWorker{
		ReposDir:       config.GIT_DIR,
		HTTPPathPrefix: config.GIT_HTTP_PATH,
		Port:           int(httpPort),
	}

	repoWorkerUpdate, repoWorkerDeleteOnly := &workers.GitRepoWorker{
		ReposDir:    config.GIT_DIR,
		DeleteOnly:  false,
		RedisClient: redisClient,
		RedisPrefix: config.REDIS_CHANNEL_PREFIX,
		RedisSuffix: config.REDIS_CHANNEL_MODULE_REGISTERED,
	}, &workers.GitRepoWorker{
		ReposDir:    config.GIT_DIR,
		DeleteOnly:  true,
		RedisClient: redisClient,
		RedisPrefix: config.REDIS_CHANNEL_PREFIX,
		RedisSuffix: config.REDIS_CHANNEL_MODULE_UNREGISTERED,
	}

	httpWorkerErrors, repoWorkerUpdateErrors, repoWorkerDeleteOnlyErrors := make(chan error, 0), make(chan error, 0), make(chan error, 0)

	httpWorkerEvents, repoWorkerUpdateEvents, repoWorkerDeleteOnlyEvents := make(chan utils.Event, 0), make(chan utils.Event, 0), make(chan utils.Event, 0)

	go httpWorker.Start(httpWorkerErrors, httpWorkerEvents)
	go repoWorkerUpdate.Start(repoWorkerUpdateErrors, repoWorkerUpdateEvents)
	go repoWorkerDeleteOnly.Start(repoWorkerDeleteOnlyErrors, repoWorkerDeleteOnlyEvents)

	for {
		select {
		case err := <-httpWorkerErrors:
			panic(err)
		case err := <-repoWorkerUpdateErrors:
			panic(err)
		case err := <-repoWorkerDeleteOnlyErrors:
			panic(err)

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
				log.Info("Started", rz.String("System", "GitRepoWorker"), rz.String("EventMessage", event.Message), rz.String("ReposDir", repoWorkerUpdate.ReposDir), rz.Bool("DeleteOnly", repoWorkerUpdate.DeleteOnly), rz.String("RedisPrefix", repoWorkerUpdate.RedisPrefix), rz.String("RedisSuffix", repoWorkerUpdate.RedisSuffix))
			}
		case event := <-repoWorkerDeleteOnlyEvents:
			switch event.Code {
			case 0:
				log.Info("Started", rz.String("System", "GitRepoWorker"), rz.String("EventMessage", event.Message), rz.String("ReposDir", repoWorkerDeleteOnly.ReposDir), rz.Bool("DeleteOnly", repoWorkerDeleteOnly.DeleteOnly), rz.String("RedisPrefix", repoWorkerDeleteOnly.RedisPrefix), rz.String("RedisSuffix", repoWorkerDeleteOnly.RedisSuffix))
			}
		}
	}
}

// StartDirectoryManagementWorker starts a new directory management worker
func StartDirectoryManagementWorker(wg *sync.WaitGroup, r *redis.Client, prefix, channel, baseDir string, deleteOnly bool) error {
	err, c, p := utils.GetRedisChannel(r, prefix, channel)
	defer p.Close()
	if err != nil {
		return err
	}

	if deleteOnly {
		log.Info("Starting directory deletion worker ...")
	} else {
		log.Info("Starting directory update worker ...")
	}

	for m := range c {
		var innerWg sync.WaitGroup

		go func(wg *sync.WaitGroup, msg *redis.Message) {
			wg.Add(1)

			n, t := utils.ParseModuleFromMessage(msg.Payload)
			if deleteOnly {
				log.Info("Deleting directory", rz.String("moduleName", n), rz.String("eventTimestamp", t))
			} else {
				log.Info("Updating directory", rz.String("moduleName", n), rz.String("eventTimestamp", t))
			}

			path := utils.GetPathForModule(baseDir, n)

			if !deleteOnly {
				err = os.RemoveAll(path)
				if err != nil {
					panic(err)
				}

				err = os.MkdirAll(path, 0777)
				if err != nil {
					panic(err)
				}

				_, err := git.PlainInit(path, false)
				if err != nil {
					panic(err)
				}
			}

			defer wg.Done()
		}(&innerWg, m)
	}

	defer wg.Done()

	return nil
}
