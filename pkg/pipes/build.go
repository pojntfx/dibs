package pipes

import (
	"os"
	"os/exec"
	"strings"
)

type Build struct {
	Command string
	Tag     string
	Context string
	File    string
}

var DefaultEnvVariables = map[string]string{
	"DOCKER_CLI_EXPERIMENTAL": "enabled",
	"DOCKER_BUILDKIT":         "1",
}

const DockerCommand = "docker"

func (build *Build) exec(args ...string) error {
	for key, value := range DefaultEnvVariables {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	command := exec.Command(args[0], args[1:]...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func (build *Build) execDocker(args ...string) error {
	return build.exec(append([]string{DockerCommand}, args...)...)
}

func (build *Build) execString(command string) error {
	return build.exec(strings.Split(command, " ")...)
}

func (build *Build) Start() error {
	return build.execString(build.Command)
}

func (build *Build) BuildImage(platform string) error {
	return build.execDocker("buildx", "build", "--progress", "plain", "--pull", "--load", "--platform", platform, "-t", build.Tag, "-f", build.File, build.Context)
}

func (build *Build) StartImage(platform string) error {
	return build.execDocker("run", "--platform", platform, "--privileged", "-v", "/var/run/docker.sock:/var/run/docker.sock", build.Tag)
}

func (build *Build) PushImage() error {
	return build.execDocker("push", build.Tag)
}

func (build *Build) CleanImage() error {
	return build.execDocker("rmi", "-f", build.Tag)
}