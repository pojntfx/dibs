package utils

import (
	"github.com/google/uuid"
	"github.com/magefile/mage/sh"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BuildConfig struct {
	Dockerfile                   string
	Architecture                 string
	Tag                          string
	BinaryInContainerPath        string
	BinaryDistPath               string
	IntegrationTestCommandBinary string
	IntegrationTestCommandDocker string
}

type BuildConfigCollection struct {
	Tag                    string
	UnitTestCommand        string
	IntegrationTestCommand string
	CleanGlob              string
	BuildConfigs           []BuildConfig
}

func (buildConfig *BuildConfig) BuildDockerImage() error {
	return sh.RunV("docker", "build", "-f", buildConfig.Dockerfile, "-t", buildConfig.Tag, ".")
}

func (buildConfig *BuildConfig) PushDockerImage() error {
	return sh.RunV("docker", "push", buildConfig.Tag)
}

func (buildConfig *BuildConfig) GetBinaryFromDockerContainer() error {
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
		if err := sh.RunV("docker", "run", "--name", id, buildConfig.Tag, "echo"); err != nil {
			return err
		}

		if err := sh.RunV("docker", "cp", id+":"+buildConfig.BinaryInContainerPath, buildConfig.BinaryDistPath); err != nil {
			return err
		}

		if err := sh.RunV("docker", "rm", "-f", id); err != nil {
			return err
		}
	} else {
		return buildConfig.GetBinaryFromDockerContainer()
	}

	return nil
}

func (buildConfig *BuildConfig) IntegrationTestBinary() error {
	cmds := strings.Split(buildConfig.IntegrationTestCommandBinary, " ")

	return sh.RunV(cmds[0], cmds[1:]...)
}

func (buildConfig *BuildConfig) IntegrationTestDocker() error {
	cmds := strings.Split(buildConfig.IntegrationTestCommandDocker, " ")

	return sh.RunV(cmds[0], cmds[1:]...)
}

func (buildConfigCollection *BuildConfigCollection) UnitTest() error {
	cmds := strings.Split(buildConfigCollection.UnitTestCommand, " ")

	return sh.RunV(cmds[0], cmds[1:]...)
}

func (buildConfigCollection *BuildConfigCollection) IntegrationTest() error {
	cmds := strings.Split(buildConfigCollection.IntegrationTestCommand, " ")

	return sh.RunV(cmds[0], cmds[1:]...)
}

func (buildConfigCollection *BuildConfigCollection) SetupMultiArch() error {
	return sh.RunV("docker", "run", "--rm", "--privileged", "multiarch/qemu-user-static", "--reset", "-p", "yes")
}

func (buildConfigCollection *BuildConfigCollection) Clean() error {
	filesToRemove, _ := filepath.Glob(buildConfigCollection.CleanGlob)
	for _, fileToRemove := range filesToRemove {
		if err := os.Remove(fileToRemove); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) BuildAllDockerImages() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.BuildDockerImage(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) PushAllDockerImages() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.PushDockerImage(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) IntegrationTestDockerAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.IntegrationTestDocker(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) IntegrationTestBinaryAll() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.IntegrationTestBinary(); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) BuildDockerManifest() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := sh.RunWith(map[string]string{
			"DOCKER_CLI_EXPERIMENTAL": "enabled",
		}, "docker", "manifest", "create", "--amend", buildConfigCollection.Tag, buildConfig.Tag); err != nil {
			return err
		}
	}

	return nil
}

func (buildConfigCollection *BuildConfigCollection) PushDockerManifest() error {
	return sh.RunV("docker", "manifest", "push", buildConfigCollection.Tag)
}

func (buildConfigCollection *BuildConfigCollection) GetAllBinariesFromDockerContainers() error {
	for _, buildConfig := range buildConfigCollection.BuildConfigs {
		if err := buildConfig.GetBinaryFromDockerContainer(); err != nil {
			return err
		}
	}

	return nil
}
