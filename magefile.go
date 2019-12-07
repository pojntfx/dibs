//+build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/pojntfx/godibs/pkg/utils"
	"os"
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
	PLATFORM     = os.Getenv("PLATFORM")
	ARCHITECTURE = os.Getenv("ARCHITECTURE")
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
