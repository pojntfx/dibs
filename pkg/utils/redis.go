package utils

import redis "github.com/go-redis/redis/v7"

// NewRedisClient returns a new Redis client
func NewRedisClient(addr string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return redisClient
}

// GetRedisChannel gets a new Go channel for a redis prefix and channel
func GetRedisChannel(r *redis.Client, prefix, channel string) (error, <-chan *redis.Message, *redis.PubSub) {
	p := r.Subscribe(prefix + ":" + channel)

	if _, err := p.Receive(); err != nil {

		return err, nil, p
	}

	return nil, p.Channel(), p
}
