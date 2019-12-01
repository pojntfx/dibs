package workers

import (
	redis "github.com/go-redis/redis/v7"
	"github.com/pojntfx/godibs/pkg/utils"
	git "gopkg.in/src-d/go-git.v4"
	"os"
	"sync"
)

// GitRepoWorker creates, updates and deletes Git repos
type GitRepoWorker struct {
	ReposDir    string // Directory in which the managed repos should reside
	DeleteOnly  bool   // Whether the worker should only perform delete operations
	Redis       utils.Redis
	RedisSuffix string // Channel suffix to use to listen to messages
}

// Start starts a GitRepoWorker
func (worker *GitRepoWorker) Start(errors chan error, events chan utils.Event) {
	events <- utils.Event{
		Code:    0,
		Message: "Started",
	}

	err, c, p := worker.Redis.GetRedisChannel(worker.RedisSuffix)
	if err != nil {
		errors <- err
	}
	defer p.Close()

	for m := range c {
		var innerWg sync.WaitGroup

		go func(wg *sync.WaitGroup, msg *redis.Message) {
			wg.Add(1)

			moduleName, moduleTimestamp := utils.ParseModuleFromMessage(msg.Payload)
			if worker.DeleteOnly {
				events <- utils.Event{
					Code:    1,
					Message: "Deleting directory for module " + moduleName + " (version " + moduleTimestamp + ")",
				}
			} else {
				events <- utils.Event{
					Code:    1,
					Message: "Updating directory for module " + moduleName + " (version " + moduleTimestamp + ")",
				}
			}

			path := utils.GetPathForModule(worker.ReposDir, moduleName)

			if !worker.DeleteOnly {
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

	events <- utils.Event{
		Code:    2,
		Message: "Server stopped",
	}
}
