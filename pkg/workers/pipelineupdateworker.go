package workers

import (
	"fmt"
	"github.com/pojntfx/godibs/pkg/utils"
)

// PipelineUpdateWorker runs an event on every update
type PipelineUpdateWorker struct {
	Pipeline    utils.Pipeline // Pipeline to run on redis event
	Redis       utils.Redis    // Redis instance to get the channel from
	RedisSuffix string         // Redis channel suffix to use to listen to messages
}

// Start starts a PipelineUpdateWorker
func (worker *PipelineUpdateWorker) Start(errors chan error, events chan utils.Event) {
	events <- utils.Event{
		Code:    0,
		Message: "Started",
	}

	err, channel, pubSub := worker.Redis.GetRedisChannel(worker.RedisSuffix)
	if err != nil {
		errors <- err
	}
	defer pubSub.Close()

	for message := range channel {
		fmt.Println(message)
	}

	events <- utils.Event{
		Code:    2,
		Message: "Pipeline update worker stopped",
	}
}
