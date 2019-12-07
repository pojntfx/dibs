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
	Dockerfile:                   "Dockerfile.amd64",
	Architecture:                 "amd64",
	Tag:                          "pojntfx/godibs:amd64",
	BinaryInContainerPath:        "/usr/local/bin/godibs",
	BinaryDistPath:               filepath.Join(".bin", "godibs-amd64"),
	IntegrationTestCommandBinary: ".bin/godibs-amd64 --help",
	IntegrationTestCommandDocker: "docker run pojntfx/godibs:amd64",
}

var buildConfigARM64 = BuildConfig{
	Dockerfile:                   "Dockerfile.arm64",
	Architecture:                 "arm64",
	Tag:                          "pojntfx/godibs:arm64",
	BinaryInContainerPath:        "/usr/local/bin/godibs",
	BinaryDistPath:               filepath.Join(".bin", "godibs-arm64"),
	IntegrationTestCommandBinary: ".bin/godibs-arm64 --help",
	IntegrationTestCommandDocker: "docker run pojntfx/godibs:arm64",
}

var buildConfigARM = BuildConfig{
	Dockerfile:                   "Dockerfile.arm",
	Architecture:                 "arm",
	Tag:                          "pojntfx/godibs:arm",
	BinaryInContainerPath:        "/usr/local/bin/godibs",
	BinaryDistPath:               filepath.Join(".bin", "godibs-arm"),
	IntegrationTestCommandBinary: ".bin/godibs-arm --help",
	IntegrationTestCommandDocker: "docker run pojntfx/godibs:arm",
}

var buildConfigCollection = BuildConfigCollection{
	Tag:                    "pojntfx/godibs",
	UnitTestCommand:        "go test ./...",
	IntegrationTestCommand: "go run main.go server --help",
	CleanGlob:              ".bin",
	BuildConfigs: []BuildConfig{
		buildConfigAMD64,
		buildConfigARM64,
		buildConfigARM,
	},
}

func BuildDockerImageAMD64() error {
	return buildConfigAMD64.BuildDockerImage()
}

func BuildDockerImageARM64() error {
	return buildConfigARM64.BuildDockerImage()
}

func BuildDockerImageARM() error {
	return buildConfigARM.BuildDockerImage()
}

func PushDockerImageAMD64() error {
	return buildConfigAMD64.PushDockerImage()
}

func PushDockerImageARM64() error {
	return buildConfigARM64.PushDockerImage()
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

func IntegrationTestDockerAMD64() error {
	return buildConfigAMD64.IntegrationTestDocker()
}

func IntegrationTestDockerARM64() error {
	return buildConfigARM64.IntegrationTestDocker()
}

func IntegrationTestDockerARM() error {
	return buildConfigARM.IntegrationTestDocker()
}

func IntegrationTestBinaryAMD64() error {
	return buildConfigAMD64.IntegrationTestDocker()
}

func IntegrationTestBinaryARM64() error {
	return buildConfigARM64.IntegrationTestDocker()
}

func IntegrationTestBinaryARM() error {
	return buildConfigARM.IntegrationTestDocker()
}

func UnitTest() error {
	return buildConfigCollection.UnitTest()
}

func IntegrationTest() error {
	return buildConfigCollection.IntegrationTest()
}

func IntegrationTestDockerAll() error {
	return buildConfigCollection.IntegrationTestDockerAll()
}

func IntegrationTestBinaryAll() error {
	return buildConfigCollection.IntegrationTestBinaryAll()
}

func SetupMultiArch() error {
	return buildConfigCollection.SetupMultiArch()
}

func Clean() error {
	return buildConfigCollection.Clean()
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
