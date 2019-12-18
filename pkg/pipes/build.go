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
	CommandDocker   = "docker"
	CommandHelm     = "helm"
	CommandSkaffold = "skaffold"
	CommandGHR      = "ghr"
	CommandCR       = "cr"

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

func (build *Build) execStdoutStderr(platform string, args ...string) error {
	for key, value := range DefaultEnvVariables {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	if err := os.Setenv(TargetPlatformEnvKey, platform); err != nil {
		return err
	}

	command := exec.Command(args[0], args[1:]...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func (build *Build) execDocker(platform string, args ...string) (string, error) {
	return build.exec(platform, append([]string{CommandDocker}, args...)...)
}

func (build *Build) execHelm(platform string, args ...string) (string, error) {
	return build.exec(platform, append([]string{CommandHelm}, args...)...)
}

func (build *Build) execSkaffold(platform string, args ...string) error {
	return build.execStdoutStderr(platform, append([]string{CommandSkaffold}, args...)...)
}

func (build *Build) execGHR(platform string, args ...string) (string, error) {
	return build.exec(platform, append([]string{CommandGHR}, args...)...)
}

func (build *Build) execCR(platform string, args ...string) (string, error) {
	return build.exec(platform, append([]string{CommandCR}, args...)...)
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

func (build *Build) StartImage(platform string, envVariableKeyValues ...struct {
	Key   string
	Value string
}) (string, error) {
	if envVariableKeyValues != nil {
		var envVariableFlags []string

		envVariableFlags = append(envVariableFlags, "-e", TargetPlatformEnvKey+"="+platform)
		for _, keyVal := range envVariableKeyValues {
			envVariableFlags = append(envVariableFlags, "-e", keyVal.Key+"="+keyVal.Value)
		}

		return build.execDocker(platform, append(append([]string{"run", "--platform", platform}, envVariableFlags...), "--privileged", "-v", "/var/run/docker.sock:/var/run/docker.sock", build.Tag)...)
	}

	return build.execDocker(platform, "run", "--platform", platform, "-e", TargetPlatformEnvKey+"="+platform, "--privileged", "-v", "/var/run/docker.sock:/var/run/docker.sock", build.Tag)
}

func (build *Build) StartChart(platform, profile string) error {
	return build.execSkaffold(platform, "run", "--profile", profile)
}

func (build *Build) DevChart(platform, profile string) error {
	return build.execSkaffold(platform, "dev", "--profile="+profile, "--port-forward="+"true", "--no-prune-children="+"true")
}

func (build *Build) PushImage(platform string) (string, error) {
	return build.execDocker(platform, "push", build.Tag)
}

func (build *Build) CleanImage(platform string) (string, error) {
	return build.execDocker(platform, "rmi", "-f", build.Tag)
}

func (build *Build) CleanStartedChart(platform, profile string) error {
	return build.execSkaffold(platform, "delete", "--profile", profile)
}
