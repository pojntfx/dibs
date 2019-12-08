package utils

import (
	"os"
	"os/exec"
	"strings"
)

type BuildConfigV2 struct {
	Platform string

	BuildCommand       string
	BuildDockerfile    string
	BuildDockerContext string

	TestUnitCommand       string
	TestUnitDockerfile    string
	TestUnitDockerContext string

	TestIntegrationGoCommand       string
	TestIntegrationGoDockerContext string
	TestIntegrationGoDockerfile    string

	TestIntegrationDockerCommand       string
	TestIntegrationDockerDockerContext string
	TestIntegrationDockerDockerfile    string

	TestIntegrationBinaryCommand       string
	TestIntegrationBinaryDockerContext string
	TestIntegrationBinaryDockerfile    string
}

func (buildConfig *BuildConfigV2) exec(commands ...string) error {
	command := exec.Command(commands[0], commands[1:]...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func (buildConfig *BuildConfigV2) execString(command string) error {
	commands := strings.Split(command, " ")

	return buildConfig.exec(commands...)
}

func (buildConfig *BuildConfigV2) execDocker(args ...string) error {
	os.Setenv("DOCKER_CLI_EXPERIMENTAL", "enabled")
	os.Setenv("DOCKER_BUILDKIT", "1")

	return buildConfig.exec(append([]string{"docker"}, args...)...)
}

func (buildConfig *BuildConfigV2) Build() error {
	return buildConfig.execString(buildConfig.BuildCommand)
}

func (buildConfig *BuildConfigV2) BuildInDocker() error {
	return buildConfig.execDocker("build", "--progress", "plain", "--pull", "--platform", buildConfig.Platform, "-f", buildConfig.BuildDockerfile, buildConfig.BuildDockerContext)
}

func (buildConfig *BuildConfigV2) TestUnit() error {
	return buildConfig.execString(buildConfig.TestUnitCommand)
}

func (buildConfig *BuildConfigV2) TestUnitInDocker() error {
	return buildConfig.execDocker("build", "--progress", "plain", "--pull", "--platform", buildConfig.Platform, "-f", buildConfig.TestUnitDockerfile, buildConfig.TestUnitDockerContext)
}

func (buildConfig *BuildConfigV2) TestIntegrationGo() error {
	return buildConfig.execString(buildConfig.TestIntegrationGoCommand)
}

func (buildConfig *BuildConfigV2) TestIntegrationGoInDocker() error {
	return buildConfig.execDocker("build", "--progress", "plain", "--pull", "--platform", buildConfig.Platform, "-f", buildConfig.TestIntegrationGoDockerfile, buildConfig.TestIntegrationGoDockerContext)
}
