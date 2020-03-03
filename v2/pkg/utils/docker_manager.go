package utils

import (
	"errors"
	"os"
)

// DockerManager manages Docker
type DockerManager struct {
	dir                    string
	stdoutChan, stderrChan chan string
}

// NewDockerManager creates a new DockerManager
func NewDockerManager(dir string, stdoutChan, stderrChan chan string) *DockerManager {
	return &DockerManager{
		dir:        dir,
		stdoutChan: stdoutChan,
		stderrChan: stderrChan,
	}
}

// Build builds and tags a Docker image
func (d *DockerManager) Build(file, context, tag string) error {
	command := NewManageableCommand("docker build -f "+file+" -t "+tag+" "+context, d.dir, d.stdoutChan, d.stderrChan)

	if err := command.Start(); err != nil {
		return err
	}

	return command.Wait()
}

// Push pushes a Docker image
func (d *DockerManager) Push(tag string) error {
	command := NewManageableCommand("docker push "+tag, d.dir, d.stdoutChan, d.stderrChan)

	if err := command.Start(); err != nil {
		return err
	}

	return command.Wait()
}

// Run runs a command in a Docker image
func (d *DockerManager) Run(tag, execLine string, dockerInDocker bool) error {
	command := NewManageableCommand("docker run -e TARGETPLATFORM="+os.Getenv("TARGETPLATFORM")+" "+tag+" "+execLine, d.dir, d.stdoutChan, d.stderrChan)
	// TODO: Add test for Docker in Docker run
	if dockerInDocker {
		command = NewManageableCommand("docker run --privileged -v /var/run/docker.sock:/var/run/docker.sock -e TARGETPLATFORM="+os.Getenv("TARGETPLATFORM")+" "+tag+" "+execLine, d.dir, d.stdoutChan, d.stderrChan)
	}

	if err := command.Start(); err != nil {
		return err
	}

	return command.Wait()
}

// CopyFromImage copies an asset from a Docker image
func (d *DockerManager) CopyFromImage(tag, assetInImage, assetOut string) error {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	runInBackgroundCommand := NewManageableCommand("docker run -d "+tag+" "+"ls", d.dir, stdoutChan, stderrChan)

	if err := runInBackgroundCommand.Start(); err != nil {
		return err
	}

	containerId := func() string {
		for {
			select {
			case id := <-stdoutChan:
				return id
			}
		}
	}()

	if err := runInBackgroundCommand.Wait(); err != nil {
		return err
	}

	if containerId == "" {
		return errors.New("could not get ID from running the image")
	}

	copyCommand := NewManageableCommand("docker cp "+containerId+":"+assetInImage+" "+assetOut, d.dir, d.stdoutChan, d.stderrChan)

	if err := copyCommand.Start(); err != nil {
		return err
	}

	return copyCommand.Wait()
}
