package utils

import (
	"github.com/google/uuid"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BuildConfig struct {
	Tag      string
	Platform string

	BinaryInContainerPath string
	BinaryDistPath        string
	CleanGlob             string

	BuildCommand       string
	BuildDockerfile    string
	BuildDockerContext string

	BuildImageCommand string

	TestUnitCommand       string
	TestUnitDockerfile    string
	TestUnitDockerContext string
	TestUnitImageTag      string

	TestIntegrationLangCommand       string
	TestIntegrationLangDockerContext string
	TestIntegrationLangDockerfile    string
	TestIntegrationLangImageTag      string

	TestIntegrationImageCommand       string
	TestIntegrationImageDockerContext string
	TestIntegrationImageDockerfile    string
	TestIntegrationImageImageTag      string

	TestIntegrationBinaryCommand       string
	TestIntegrationBinaryDockerContext string
	TestIntegrationBinaryDockerfile    string
	TestIntegrationBinaryImageTag      string
}

type BuildConfigCollection struct {
	ManifestTag string

	BuildConfigs []BuildConfig
}

func (buildConfig *BuildConfig) exec(commands ...string) error {
	os.Setenv("DOCKER_CLI_EXPERIMENTAL", "enabled")
	os.Setenv("DOCKER_BUILDKIT", "1")
	os.Setenv("TARGETPLATFORM", buildConfig.Platform)

	command := exec.Command(commands[0], commands[1:]...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func (buildConfig *BuildConfig) execString(command string) error {
	commands := strings.Split(command, " ")

	return buildConfig.exec(commands...)
}

func (buildConfig *BuildConfig) execDocker(args ...string) error {
	return buildConfig.exec(append([]string{"docker"}, args...)...)
}

func (buildConfig *BuildConfig) Build() error {
	return buildConfig.execString(buildConfig.BuildCommand)
}

func (buildConfig *BuildConfig) BuildInDocker() error {
	return buildConfig.execDocker("buildx", "build", "--progress", "plain", "--pull", "--load", "--platform", buildConfig.Platform, "-t", buildConfig.Tag, "-f", buildConfig.BuildDockerfile, buildConfig.BuildDockerContext)
}

func (buildConfig *BuildConfig) BuildImage() error {
	return buildConfig.execString(buildConfig.BuildImageCommand)
}

func (buildConfig *BuildConfig) TestUnit() error {
	return buildConfig.execString(buildConfig.TestUnitCommand)
}

func (buildConfig *BuildConfig) TestUnitInDocker() error {
	if err := buildConfig.execDocker("buildx", "build", "--progress", "plain", "--pull", "--load", "--platform", buildConfig.Platform, "-t", buildConfig.TestUnitImageTag, "-f", buildConfig.TestUnitDockerfile, buildConfig.TestUnitDockerContext); err != nil {
		return err
	}

	return buildConfig.execDocker("run", "--platform", buildConfig.Platform, "-e", "TARGETPLATFORM="+buildConfig.Platform, "--privileged", "-v", "/var/run/docker.sock:/var/run/docker.sock", buildConfig.TestUnitImageTag)
}

func (buildConfig *BuildConfig) TestIntegrationGo() error {
	return buildConfig.execString(buildConfig.TestIntegrationLangCommand)
}

func (buildConfig *BuildConfig) TestIntegrationGoInDocker() error {
	if err := buildConfig.execDocker("buildx", "build", "--progress", "plain", "--pull", "--load", "--platform", buildConfig.Platform, "-t", buildConfig.TestIntegrationLangImageTag, "-f", buildConfig.TestIntegrationLangDockerfile, buildConfig.TestIntegrationLangDockerContext); err != nil {
		return err
	}

	return buildConfig.execDocker("run", "--platform", buildConfig.Platform, "-e", "TARGETPLATFORM="+buildConfig.Platform, "--privileged", "-v", "/var/run/docker.sock:/var/run/docker.sock", buildConfig.TestIntegrationLangImageTag)
}

func (buildConfig *BuildConfig) TestIntegrationBinary() error {
	return buildConfig.execString(buildConfig.TestIntegrationBinaryCommand)
}

func (buildConfig *BuildConfig) TestIntegrationBinaryInDocker() error {
	if err := buildConfig.execDocker("buildx", "build", "--progress", "plain", "--pull", "--load", "--platform", buildConfig.Platform, "-t", buildConfig.TestIntegrationBinaryImageTag, "-f", buildConfig.TestIntegrationBinaryDockerfile, buildConfig.TestIntegrationBinaryDockerContext); err != nil {
		return err
	}

	return buildConfig.execDocker("run", "--platform", buildConfig.Platform, "-e", "TARGETPLATFORM="+buildConfig.Platform, "--privileged", "-v", "/var/run/docker.sock:/var/run/docker.sock", buildConfig.TestIntegrationBinaryImageTag)
}

func (buildConfig *BuildConfig) TestIntegrationImage() error {
	return buildConfig.execString(buildConfig.TestIntegrationImageCommand)
}

func (buildConfig *BuildConfig) TestIntegrationImageInDocker() error {
	if err := buildConfig.execDocker("buildx", "build", "--progress", "plain", "--pull", "--load", "--platform", buildConfig.Platform, "-t", buildConfig.TestIntegrationImageImageTag, "-f", buildConfig.TestIntegrationImageDockerfile, buildConfig.TestIntegrationImageDockerContext); err != nil {
		return err
	}

	return buildConfig.execDocker("run", "--platform", buildConfig.Platform, "-e", "TARGETPLATFORM="+buildConfig.Platform, "--privileged", "-v", "/var/run/docker.sock:/var/run/docker.sock", buildConfig.TestIntegrationImageImageTag)
}

func (buildConfig *BuildConfig) PushDockerImage() error {
	return buildConfig.execDocker("push", buildConfig.Tag)
}

func (buildConfig *BuildConfig) GetBinaryFromDockerImage() error {
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
		if err := buildConfig.execDocker("run", "--platform", buildConfig.Platform, "-e", "TARGETPLATFORM="+buildConfig.Platform, "--name", id, buildConfig.Tag, "echo"); err != nil {
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

func (buildConfig *BuildConfig) Clean() error {
	filesToRemove, _ := filepath.Glob(buildConfig.CleanGlob)

	for _, fileToRemove := range filesToRemove {
		if err := os.Remove(fileToRemove); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) getBuildConfigForArchitecture(architecture string) BuildConfig {
	var buildConfigForArchitecture BuildConfig

	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if buildConfig.Platform == architecture {
			buildConfigForArchitecture = buildConfig
			break
		}
	}

	return buildConfigForArchitecture
}

func (buildConfigCollection *BuildConfigCollection) Build(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.Build()
}

func (buildConfigCollection *BuildConfigCollection) BuildInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.BuildInDocker()
}

func (buildConfigCollection *BuildConfigCollection) BuildImage(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.BuildImage()
}

func (buildConfigCollection *BuildConfigCollection) TestUnit(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestUnit()
}

func (buildConfigCollection *BuildConfigCollection) TestUnitInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestUnitInDocker()
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationLang(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationGo()
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationLangInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationGoInDocker()
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationBinary(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationBinary()
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationBinaryInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationBinaryInDocker()
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationImage(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationImage()
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationImageInDocker(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.TestIntegrationImageInDocker()
}

func (buildConfigCollection *BuildConfigCollection) PushDockerImage(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.PushDockerImage()
}

func (buildConfigCollection *BuildConfigCollection) GetBinaryFromDockerImage(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.GetBinaryFromDockerImage()
}

func (buildConfigCollection *BuildConfigCollection) Clean(architecture string) error {
	buildConfig := buildConfigCollection.getBuildConfigForArchitecture(architecture)

	return buildConfig.Clean()
}

func (buildConfigCollection *BuildConfigCollection) BuildDockerManifest() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.execDocker("manifest", "create", "--amend", buildConfigCollection.ManifestTag, buildConfig.Tag); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) PushDockerManifest() error {
	return buildConfigCollection.BuildConfigs[0].execDocker("manifest", "push", buildConfigCollection.ManifestTag)
}

func (buildConfigCollection *BuildConfigCollection) BuildAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.Build(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) BuildInDockerAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.BuildInDocker(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) BuildImageAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.BuildImage(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) TestUnitAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.TestUnit(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) TestUnitInDockerAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.TestUnitInDocker(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationLangAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.TestIntegrationGo(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationLangInDockerAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.TestIntegrationGoInDocker(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationBinaryAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.TestIntegrationBinary(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationBinaryInDockerAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.TestIntegrationBinaryInDocker(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationImageAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.TestIntegrationImage(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) TestIntegrationImageInDockerAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.TestIntegrationImageInDocker(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) PushDockerImageAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.PushDockerImage(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) GetBinaryFromDockerImageAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.GetBinaryFromDockerImage(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) CleanAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.Clean(); err != nil {
			return err
		}
	}

	return nil
}
