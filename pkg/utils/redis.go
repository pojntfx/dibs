package utils

import redis "github.com/go-redis/redis/v7"

type Redis struct {
	client *redis.Client
	Addr   string
	Prefix string
}

// Connect creates a new Redis client
func (r *Redis) Connect() {
	r.client = redis.NewClient(&redis.Options{
		Addr: r.Addr,
	})
}

// GetRedisChannel gets a new Go channel for a Redis channel
func (r *Redis) GetRedisChannel(channel string) (error, <-chan *redis.Message, *redis.PubSub) {
	pubSub := r.client.Subscribe(r.Prefix + ":" + channel)

	if _, err := pubSub.Receive(); err != nil {

		return err, nil, pubSub
	}

	return nil, pubSub.Channel(), pubSub
}

// getChannelWithPrefix returns a channel suffix with the prefix
func (r *Redis) getChannelWithPrefix(suffix string) string {
	return r.Prefix + ":" + suffix
}

// getMessageWithTimestamp returns a message with the corresponding timestamp
func (r *Redis) getMessageWithTimestamp(message string) string {
	return WithTimestamp(message)
}

// PublishWithTimestamp publishes a message to a channel and adds a timestamp to it
func (r *Redis) PublishWithTimestamp(suffix, message string) {
	r.client.Publish(r.getChannelWithPrefix(suffix), r.getMessageWithTimestamp(message))
}
