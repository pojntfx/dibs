package pipes

import (
	"os"
	"os/exec"
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

const (
	DockerCommand        = "docker"
	HelmCommand          = "helm"
	TargetPlatformEnvKey = "TARGETPLATFORM"
)

func (build *Build) exec(platform string, args ...string) (string, error) {
	for key, value := range DefaultEnvVariables {
		if err := os.Setenv(key, value); err != nil {
			return "", err
		}
	}

	if err := os.Setenv(TargetPlatformEnvKey, platform); err != nil {
		return "", err
	}

	command := exec.Command(args[0], args[1:]...)

	output, err := command.CombinedOutput()

	return string(output), err
}

func (build *Build) execDocker(platform string, args ...string) (string, error) {
	return build.exec(platform, append([]string{DockerCommand}, args...)...)
}

func (build *Build) execHelm(platform string, args ...string) (string, error) {
	return build.exec(platform, append([]string{HelmCommand}, args...)...)
}

func (build *Build) execString(platform string, command string) (string, error) {
	return build.exec(platform, "sh", "-c", command)
}

func (build *Build) Start(platform string) (string, error) {
	return build.execString(platform, build.Command)
}

func (build *Build) BuildImage(platform string) (string, error) {
	return build.execDocker(platform, "buildx", "build", "--progress", "plain", "--pull", "--load", "--platform", platform, "-t", build.Tag, "-f", build.File, build.Context)
}

func (build *Build) StartImage(platform string) (string, error) {
	return build.execDocker(platform, "run", "--platform", platform, "-e", "TARGETPLATFORM="+platform, "--privileged", "-v", "/var/run/docker.sock:/var/run/docker.sock", build.Tag)
}

func (build *Build) PushImage(platform string) (string, error) {
	return build.execDocker(platform, "push", build.Tag)
}

func (build *Build) CleanImage(platform string) (string, error) {
	return build.execDocker(platform, "rmi", "-f", build.Tag)
}
