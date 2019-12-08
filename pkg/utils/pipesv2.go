package utils

import (
	"github.com/google/uuid"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BuildConfigV2 struct {
	Tag      string
	Platform string

	BinaryInContainerPath string
	BinaryDistPath        string
	CleanGlob             string

	BuildCommand       string
	BuildDockerfile    string
	BuildDockerContext string

	BuildDockerTag           string
	BuildDockerCommand       string
	BuildDockerDockerContext string
	BuildDockerDockerfile    string

	TestUnitCommand       string
	TestUnitDockerfile    string
	TestUnitDockerContext string

	TestIntegrationGoCommand       string
	TestIntegrationGoDockerContext string
	TestIntegrationGoDockerfile    string

	TestIntegrationDockerCommand       string
	TestIntegrationDockerDockerContext string
	TestIntegrationDockerDockerfile    string

	TestIntegrationDockerTag           string
	TestIntegrationBinaryCommand       string
	TestIntegrationBinaryDockerContext string
	TestIntegrationBinaryDockerfile    string
}

type BuildConfigCollectionV2 struct {
	BuildConfigs []BuildConfigV2
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
	return buildConfig.execDocker("build", "--progress", "plain", "--pull", "--platform", buildConfig.Platform, "-t", buildConfig.Tag, "-f", buildConfig.BuildDockerfile, buildConfig.BuildDockerContext)
}

func (buildConfig *BuildConfigV2) BuildDocker() error {
	return buildConfig.execString(buildConfig.BuildDockerCommand)
}

func (buildConfig *BuildConfigV2) BuildDockerInDocker() error {
	if err := buildConfig.execDocker("build", "--progress", "plain", "--pull", "--platform", buildConfig.Platform, "-t", buildConfig.BuildDockerTag, "-f", buildConfig.BuildDockerDockerfile, buildConfig.BuildDockerDockerContext); err != nil {
		return err
	}

	return buildConfig.execDocker("run", "--platform", buildConfig.Platform, "--privileged", "-v", "/var/run/docker.sock:/var/run/docker.sock", buildConfig.BuildDockerTag)
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

func (buildConfig *BuildConfigV2) TestIntegrationBinary() error {
	return buildConfig.execString(buildConfig.TestIntegrationBinaryCommand)
}

func (buildConfig *BuildConfigV2) TestIntegrationBinaryInDocker() error {
	return buildConfig.execDocker("build", "--progress", "plain", "--pull", "--platform", buildConfig.Platform, "-f", buildConfig.TestIntegrationBinaryDockerfile, buildConfig.TestIntegrationBinaryDockerContext)
}

func (buildConfig *BuildConfigV2) TestIntegrationDocker() error {
	return buildConfig.execString(buildConfig.TestIntegrationDockerCommand)
}

func (buildConfig *BuildConfigV2) TestIntegrationDockerInDocker() error {
	if err := buildConfig.BuildDocker(); err != nil {
		return nil
	}

	if err := buildConfig.execDocker("build", "--progress", "plain", "--pull", "--platform", buildConfig.Platform, "-t", buildConfig.TestIntegrationDockerTag, "-f", buildConfig.TestIntegrationDockerDockerfile, buildConfig.TestIntegrationDockerDockerContext); err != nil {
		return err
	}

	return buildConfig.execDocker("run", "--platform", buildConfig.Platform, "--privileged", "-v", "/var/run/docker.sock:/var/run/docker.sock", buildConfig.TestIntegrationDockerTag)
}

func (buildConfig *BuildConfigV2) PushDockerImage() error {
	return buildConfig.execDocker("push", buildConfig.Tag)
}

func (buildConfig *BuildConfigV2) GetBinaryFromDockerImage() error {
	id := uuid.New().String()
	distDir, _ := filepath.Split(buildConfig.BinaryDistPath)
	if err := os.MkdirAll(distDir, 0777); err != nil {
		return err
	}

	out, err := exec.Command("docker", "ps", "-aqf", "name="+id).Output()
	if err != nil {
		return err
	}
	if string(out) != "\n" {
		if err := buildConfig.execDocker("run", "--platform", buildConfig.Platform, "--name", id, buildConfig.Tag, "echo"); err != nil {
			return err
		}

		if err := buildConfig.execDocker("cp", id+":"+buildConfig.BinaryInContainerPath, buildConfig.BinaryDistPath); err != nil {
			return err
		}

		if err := buildConfig.execDocker("rm", "-f", id); err != nil {
			return err
		}
	} else {
		return buildConfig.GetBinaryFromDockerImage()
	}

	return nil
}

func (buildConfig *BuildConfigV2) Clean() error {
	filesToRemove, _ := filepath.Glob(buildConfig.CleanGlob)

	for _, fileToRemove := range filesToRemove {
		if err := os.Remove(fileToRemove); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollectionV2) getBuildConfigForArchitecture(architecture string) BuildConfigV2 {
	var buildConfigForArchitecture BuildConfigV2

	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if buildConfig.Platform == architecture {
			buildConfigForArchitecture = buildConfig
			break
		}
	}

	return buildConfigForArchitecture
}

func (buildConfigCollection *BuildConfigCollectionV2) Build(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.Build()
}

func (buildConfigCollection *BuildConfigCollectionV2) BuildInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.BuildInDocker()
}

func (buildConfigCollection *BuildConfigCollectionV2) BuildDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.BuildDocker()
}

func (buildConfigCollection *BuildConfigCollectionV2) BuildDockerInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.BuildDockerInDocker()
}

func (buildConfigCollection *BuildConfigCollectionV2) TestUnit(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestUnit()
}

func (buildConfigCollection *BuildConfigCollectionV2) TestUnitInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestUnitInDocker()
}

func (buildConfigCollection *BuildConfigCollectionV2) TestIntegrationGo(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationGo()
}

func (buildConfigCollection *BuildConfigCollectionV2) TestIntegrationGoInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationGoInDocker()
}

func (buildConfigCollection *BuildConfigCollectionV2) TestIntegrationBinary(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationBinary()
}

func (buildConfigCollection *BuildConfigCollectionV2) TestIntegrationBinaryInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationBinaryInDocker()
}

func (buildConfigCollection *BuildConfigCollectionV2) TestIntegrationDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationDocker()
}

func (buildConfigCollection *BuildConfigCollectionV2) TestIntegrationDockerInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationDockerInDocker()
}

func (buildConfigCollection *BuildConfigCollectionV2) PushDockerImage(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.PushDockerImage()
}

func (buildConfigCollection *BuildConfigCollectionV2) GetBinaryFromDockerImage(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.GetBinaryFromDockerImage()
}

func (buildConfigCollection *BuildConfigCollectionV2) Clean(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.Clean()
}
