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
		name, _ := utils.ParseModuleFromMessage(message.Payload)
		worker.Pipeline.RunCommandsOnly()
		// TODO:
		// - Get the module names that are being replaced from the configuration system
		// - Check if the module in message.Payload is in the replace module names array
		// - Run the pipeline commands only (no need)
		// - Add an ignore blob to ignore the changes by the `mage build` etc. commands later on (otherwise there would be a feedback loop here, one should specify `*.go.pb` which would then be checked at the path evaluation in the FSWatcher, which is currently only the PIPELINE_UP_DIR_PUSH)
		fmt.Println(name)
	}

	events <- utils.Event{
		Code:    2,
		Message: "Pipeline update worker stopped",
	}
}
