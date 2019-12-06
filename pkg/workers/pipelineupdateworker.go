package workers

import (
	"github.com/pojntfx/godibs/pkg/utils"
	git "gopkg.in/src-d/go-git.v4"
	"os"
	"path/filepath"
)

// PipelineUpdateWorker runs an event on every update
type PipelineUpdateWorker struct {
	Modules       []string       // Modules that are to be downloaded
	Pipeline      utils.Pipeline // Pipeline to run on redis event
	LocalCloneDir string         // Directory to clone the local modules into
	HTTPBaseURL   string         // Base URL to clone the modules from
	Redis         utils.Redis    // Redis instance to get the channel from
	RedisSuffix   string         // Redis channel suffix to use to listen to messages
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
		pushedModule, _ := utils.ParseModuleFromMessage(message.Payload)
		for _, managedModule := range worker.Modules {
			events <- utils.Event{
				Code:    1,
				Message: worker.Modules[0],
			}
			if pushedModule == managedModule {
				localModulePath := filepath.Join(worker.LocalCloneDir, pushedModule)
				remoteModulePath := worker.HTTPBaseURL + "/" + pushedModule

				events <- utils.Event{
					Code:    1,
					Message: "Cloning module from " + remoteModulePath + " to " + localModulePath,
				}

				if err := os.RemoveAll(localModulePath); err != nil {
					errors <- err
				}

				if err := os.MkdirAll(localModulePath, 0777); err != nil {
					errors <- err
				}

				if _, err := git.PlainClone(localModulePath, false, &git.CloneOptions{
					URL:               remoteModulePath,
					RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
				}); err != nil {
					errors <- err
				}

				if err := worker.Pipeline.RunCommandsOnly(); err != nil {
					errors <- err
				}
			}
		}
	}

	events <- utils.Event{
		Code:    2,
		Message: "Pipeline update worker stopped",
	}
}
