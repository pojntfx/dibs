package common

import (
	redis "github.com/go-redis/redis/v7"
)

const (
	REDIS_CHANNEL_MODULE_BUILT        = "module_built"
	REDIS_CHANNEL_MODULE_TESTED       = "module_tested"
	REDIS_CHANNEL_MODULE_STARTED      = "module_started"
	REDIS_CHANNEL_MODULE_REGISTERED   = "module_registered"
	REDIS_CHANNEL_MODULE_UNREGISTERED = "module_unregistered"
	GIT_COMMIT_MESSAGE                = "module_synced"
	REDIS_CHANNEL_MODULE_PUSHED       = "module_pushed"
)

// GetNewRedisClient returns a new Redis client
func GetNewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
