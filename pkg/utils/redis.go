package utils

import (
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

// Redis is a configured Redis instance holder
type Redis struct {
	client   *redis.Client
	Addr     string
	Password string
	Prefix   string
}

// Connect creates a new Redis client
func (r *Redis) Connect() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
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

// WithTimestamp gets a message name with the current timestamp
func WithTimestamp(message string) string {
	currentTime := time.Now().UnixNano()

	return message + "@" + strconv.Itoa(int(currentTime))
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
