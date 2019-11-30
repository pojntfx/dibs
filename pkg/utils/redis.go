package utils

import redis "github.com/go-redis/redis/v7"

// GetNewRedisClient returns a new Redis client
func GetNewRedisClient(addr string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return redisClient
}
