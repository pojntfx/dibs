package utils

import (
	"os"
	"os/exec"
)

type Dibs struct {
	Manifest  Manifest
	Platforms []PlatformConfig
}

type (
	PlatformConfig struct {
		Platform string
		Image    BuildConfig
		Binary   BinaryConfig
		Tests    struct {
			Unit        BuildConfig
			Integration struct {
				Lang   BuildConfig
				Image  BuildConfig
				Binary BuildConfig
			}
		}
	}

	BinaryConfig struct {
		BuildCommand string
		PathInImage  string
		DistPath     string
		CleanGlob    string
	}

	BuildConfig struct {
		BuildCommand string
		Tag          string
		Context      string
		File         string
	}
)

var (
	DefaultEnvVariables = map[string]string{
		"DOCKER_CLI_EXPERIMENTAL": "enabled",
		"DOCKER_BUILDKIT":         "1",
	}
)

const (
	DockerCommand = "docker"
)

func (platformConfig *PlatformConfig) exec(args ...string) error {
	for key, value := range DefaultEnvVariables {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	if err := os.Setenv("TARGETPLATFORM", platformConfig.Platform); err != nil {
		return err
	}

	command := exec.Command(args[0], args[1:]...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func (platformConfig *PlatformConfig) execDocker(args ...string) error {
	return platformConfig.exec(append([]string{DockerCommand}, args...)...)
}
