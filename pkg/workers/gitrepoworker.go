package workers

import (
	redis "github.com/go-redis/redis/v7"
	"github.com/pojntfx/godibs/pkg/utils"
)

// GitRepoWorker creates, updates and deletes Git repos
type GitRepoWorker struct {
	ReposDir    string        // Directory in which the managed repos should reside
	DeleteOnly  bool          // Whether the worker should only perform delete operations
	RedisClient *redis.Client // Redis client to listen to messages with
	RedisPrefix string        // Channel prefix to use to listen to messages
	RedisSuffix string        // Channel suffix to use to listen to messages
}

// Start starts a GitRepoWorker
func (worker *GitRepoWorker) Start(errors chan error, events chan utils.Event) {
	events <- utils.Event{
		Code:    0,
		Message: "Started",
	}
}
