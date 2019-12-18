package pipes

import (
	"os"
	"os/exec"
)

// Build is a build target
type Build struct {
	Command string // Command to run to build
	Tag     string // Tag of the Docker image to build
	Context string // Docker context for the Docker image build
	File    string // Dockerfile for the Docker image build
}

// DefaultEnvVariables are the default env variables to set on every `exec` call
var DefaultEnvVariables = map[string]string{
	"DOCKER_CLI_EXPERIMENTAL": "enabled",
	"DOCKER_BUILDKIT":         "1",
}

const (
	CommandDocker        = "docker"   // The command to run Docker
	CommandHelm          = "helm"     // The command to run Helm
	CommandSkaffold      = "skaffold" // The command to run Skaffold
	CommandGHR           = "ghr"      // The command to run ghr
	CommandChartReleaser = "cr"       // The command to run Chart Releaser

	TargetPlatformEnvKey = "TARGETPLATFORM" // Env variable key for setting the target platform
)

// exec runs a command and returns the output
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

// execStdoutStderr runs a command and logs the `stdout` and `stderr` to `stdout`
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

// execDocker runs a Docker command
func (build *Build) execDocker(platform string, args ...string) (string, error) {
	return build.exec(platform, append([]string{CommandDocker}, args...)...)
}

// execHelm runs a Helm command
func (build *Build) execHelm(platform string, args ...string) (string, error) {
	return build.exec(platform, append([]string{CommandHelm}, args...)...)
}

// execSkaffold runs a Skaffold command
func (build *Build) execSkaffold(platform string, args ...string) error {
	return build.execStdoutStderr(platform, append([]string{CommandSkaffold}, args...)...)
}

// execGHR runs a ghr command
func (build *Build) execGHR(platform string, args ...string) (string, error) {
	return build.exec(platform, append([]string{CommandGHR}, args...)...)
}

// execChartReleaser runs a Chart Releaser command
func (build *Build) execChartReleaser(platform string, args ...string) (string, error) {
	return build.exec(platform, append([]string{CommandChartReleaser}, args...)...)
}

// execString runs a command with `sh`
func (build *Build) execString(platform string, command string) (string, error) {
	return build.exec(platform, "sh", "-c", command)
}

// Start starts the build
func (build *Build) Start(platform string) (string, error) {
	return build.execString(platform, build.Command)
}

// BuildImage builds the Docker image
func (build *Build) BuildImage(platform string) (string, error) {
	return build.execDocker(platform, "buildx", "build", "--progress", "plain", "--pull", "--load", "--platform", platform, "-t", build.Tag, "-f", build.File, build.Context)
}

// StartImage start the Docker image with environment variables
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

// StartChart starts the Helm chart
func (build *Build) StartChart(platform, profile string) error {
	return build.execSkaffold(platform, "run", "--profile", profile)
}

// DevChart starts a chart development session
func (build *Build) DevChart(platform, profile string) error {
	return build.execSkaffold(platform, "dev", "--profile="+profile, "--port-forward="+"true", "--no-prune-children="+"true")
}

// PushImage pushes the Docker image to a Docker registry
func (build *Build) PushImage(platform string) (string, error) {
	return build.execDocker(platform, "push", build.Tag)
}

// CleanImage removes the images locally
func (build *Build) CleanImage(platform string) (string, error) {
	return build.execDocker(platform, "rmi", "-f", build.Tag)
}

// CleanStartedChart stops and uninstalls the started chart
func (build *Build) CleanStartedChart(platform, profile string) error {
	return build.execSkaffold(platform, "delete", "--profile", profile)
}
