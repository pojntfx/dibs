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

func Start() error {
	return buildConfiguration.Start()
}

func UnitTests() error {
	return buildConfiguration.UnitTests()
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
	Tag:                   "pojntfx/godibs-amd64",
	BinaryInContainerPath: "/usr/local/bin/godibs",
	BinaryDistPath:        filepath.Join(".bin", "godibs-amd64"),
}

var buildConfigARM64 = BuildConfig{
	Architecture:          "arm64",
	Tag:                   "pojntfx/godibs-arm64",
	BinaryInContainerPath: "/usr/local/bin/godibs",
	BinaryDistPath:        filepath.Join(".bin", "godibs-arm64"),
}

func BuildDockerImageAMD64() error {
	return buildConfigAMD64.BuildDockerImage()
}

func BuildDockerImageARM64() error {
	return buildConfigARM64.BuildDockerImage()
}

func GetBinaryFromDockerContainerAMD64() error {
	return buildConfigAMD64.GetBinaryFromDockerContainer()
}

func GetBinaryFromDockerContainerARM64() error {
	return buildConfigARM64.GetBinaryFromDockerContainer()
}

type BuildConfig struct {
	Architecture          string
	Tag                   string
	BinaryInContainerPath string
	BinaryDistPath        string
}

func (buildConfig *BuildConfig) BuildDockerImage() error {
	return sh.RunV("docker", "build", "-f", "Dockerfile."+buildConfig.Architecture, "-t", buildConfig.Tag, ".")
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
