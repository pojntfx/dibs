package workers

import (
	"github.com/go-redis/redis/v7"
	"github.com/pojntfx/godibs/pkg/utils"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"sync"
)

// GitRepoWorker creates, updates and deletes Git repos
type GitRepoWorker struct {
	ReposDir    string      // Directory in which the managed repos should reside
	DeleteOnly  bool        // Whether the worker should only perform delete operations
	Redis       utils.Redis // Redis instance to get the channel from
	RedisSuffix string      // Redis channel suffix to use to listen to messages
}

// Start starts a GitRepoWorker
func (worker *GitRepoWorker) Start(errors chan error, events chan utils.Event) {
	events <- utils.Event{
		Code:    0,
		Message: "Started",
	}

	err, channel, pubSub := worker.Redis.GetRedisChannel(worker.RedisSuffix)
	if err != nil {
		errors <- err
	}
	defer func() {
		if err := pubSub.Close(); err != nil {
			errors <- err
		}
	}()

	for message := range channel {
		var innerWg sync.WaitGroup

		go func(wg *sync.WaitGroup, message *redis.Message) {
			wg.Add(1)

			module, moduleTimestamp := utils.ParseModuleFromMessage(message.Payload)
			if worker.DeleteOnly {
				events <- utils.Event{
					Code:    1,
					Message: "Deleting directory for module " + module + " (version " + moduleTimestamp + ")",
				}
			} else {
				events <- utils.Event{
					Code:    1,
					Message: "Updating directory for module " + module + " (version " + moduleTimestamp + ")",
				}
			}

			path := utils.GetPathForModule(worker.ReposDir, module)

			if !worker.DeleteOnly {
				err = os.RemoveAll(path)
				if err != nil {
					errors <- err
				}

				err = os.MkdirAll(path, 0777)
				if err != nil {
					errors <- err
				}

				_, err := git.PlainInit(path, false)
				if err != nil {
					errors <- err
				}
			}

			defer wg.Done()
		}(&innerWg, message)
	}

	events <- utils.Event{
		Code:    2,
		Message: "Server stopped",
	}
}
