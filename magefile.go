//+build mage

package main

import (
	"github.com/google/uuid"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pojntfx/godibs/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	buildConfiguration = utils.BuildConfiguration{
		ProjectName:      "godibs",
		MainFile:         "main.go",
		DockerRepoPrefix: "pojntfx",
		GoCmd:            mg.GoCmd(),
		DirTemp:          os.TempDir(),
		DirBin:           ".bin",
		DirInstall:       filepath.Join("/usr", "local", "bin", "godibs"),
		ProfilesBase: []string{
			"godibs",
			"godibs" + "-dev",
		},
		Architectures: []string{
			"amd64",
			"arm64",
		},
	}
	ARCHITECTURE = os.Getenv("ARCHITECTURE")
	PLATFORM     = os.Getenv("PLATFORM")
)

func Build() error {
	return buildConfiguration.Build()
}

func BinaryBuild() error {
	return buildConfiguration.BinaryBuild(PLATFORM, ARCHITECTURE)
}

func BinaryInstall() error {
	return buildConfiguration.BinaryInstall(PLATFORM, ARCHITECTURE)
}

func Clean() error {
	return buildConfiguration.Clean()
}

func IntegrationTests() error {
	return buildConfiguration.IntegrationTests()
}

func BinaryIntegrationTests() error {
	mg.SerialDeps(BinaryInstall)

	return buildConfiguration.BinaryIntegrationTests()
}

func DockerMultiarchSetup() error {
	return buildConfiguration.DockerMultiarchSetup()
}

func SkaffoldBuild() error {
	mg.SerialDeps(DockerMultiarchSetup)

	return buildConfiguration.SkaffoldBuild()
}

func DockerManifestBuild() error {
	return buildConfiguration.DockerManifestBuild()
}

var buildConfigAMD64 = BuildConfig{
	Architecture:          "amd64",
	Tag:                   "pojntfx/godibs:amd64",
	BinaryInContainerPath: "/usr/local/bin/godibs",
	BinaryDistPath:        filepath.Join(".bin", "godibs-amd64"),
}

var buildConfigARM64 = BuildConfig{
	Architecture:          "arm64",
	Tag:                   "pojntfx/godibs:arm64",
	BinaryInContainerPath: "/usr/local/bin/godibs",
	BinaryDistPath:        filepath.Join(".bin", "godibs-arm64"),
}

var buildConfigARM = BuildConfig{
	Architecture:          "arm",
	Tag:                   "pojntfx/godibs:arm",
	BinaryInContainerPath: "/usr/local/bin/godibs",
	BinaryDistPath:        filepath.Join(".bin", "godibs-arm"),
}

var buildConfigCollection = BuildConfigCollection{
	Tag:             "pojntfx/godibs",
	StartCommand:    "go run main.go server",
	UnitTestCommand: "go test ./...",
	BuildConfigs: []BuildConfig{
		buildConfigAMD64,
		buildConfigARM64,
		buildConfigARM,
	},
}

func BuildDockerImageAMD64() error {
	return buildConfigAMD64.BuildDockerImage()
}

func PushDockerImageAMD64() error {
	return buildConfigAMD64.PushDockerImage()
}

func BuildDockerImageARM64() error {
	return buildConfigARM64.BuildDockerImage()
}

func PushDockerImageARM64() error {
	return buildConfigARM64.PushDockerImage()
}

func BuildDockerImageARM() error {
	return buildConfigARM.BuildDockerImage()
}

func PushDockerImageARM() error {
	return buildConfigARM.PushDockerImage()
}

func GetBinaryFromDockerContainerAMD64() error {
	return buildConfigAMD64.GetBinaryFromDockerContainer()
}

func GetBinaryFromDockerContainerARM64() error {
	return buildConfigARM64.GetBinaryFromDockerContainer()
}

func GetBinaryFromDockerContainerARM() error {
	return buildConfigARM.GetBinaryFromDockerContainer()
}

func Start() error {
	return buildConfigCollection.Start()
}

func UnitTest() error {
	return buildConfigCollection.UnitTest()
}

func SetupMultiArch() error {
	return buildConfigCollection.SetupMultiArch()
}

func BuildAllDockerImages() error {
	return buildConfigCollection.BuildAllDockerImages()
}

func PushAllDockerImages() error {
	return buildConfigCollection.PushAllDockerImages()
}

func BuildDockerManifest() error {
	return buildConfigCollection.BuildDockerManifest()
}

func PushDockerManifest() error {
	return buildConfigCollection.PushDockerManifest()
}

func GetAllBinariesFromDockerContainers() error {
	return buildConfigCollection.GetAllBinariesFromDockerContainers()
}

type BuildConfig struct {
	Architecture          string
	Tag                   string
	BinaryInContainerPath string
	BinaryDistPath        string
}

type BuildConfigCollection struct {
	Tag             string
	StartCommand    string
	UnitTestCommand string
	BuildConfigs    []BuildConfig
}

func (buildConfig *BuildConfig) BuildDockerImage() error {
	return sh.RunV("docker", "build", "-f", "Dockerfile."+buildConfig.Architecture, "-t", buildConfig.Tag, ".")
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

func (buildConfigCollection *BuildConfigCollection) Start() error {
	cmds := strings.Split(buildConfigCollection.StartCommand, " ")

	return sh.RunV(cmds[0], cmds[1:]...)
}

func (buildConfigCollection *BuildConfigCollection) UnitTest() error {
	cmds := strings.Split(buildConfigCollection.UnitTestCommand, " ")

	return sh.RunV(cmds[0], cmds[1:]...)
}

func (buildConfigCollection *BuildConfigCollection) SetupMultiArch() error {
	return sh.RunV("docker", "run", "--rm", "--privileged", "multiarch/qemu-user-static", "--reset", "-p", "yes")
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
