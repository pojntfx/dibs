//+build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/pojntfx/godibs/pkg/utils"
	"os"
	"path/filepath"
)

var (
	PROJECT_NAME     = "godibs"
	MAIN_FILE        = "main.go"
	DOCKER_NAMESPACE = "pojntfx"
	COMMAND_GO       = mg.GoCmd()
	DIR_TEMP         = os.TempDir()
	DIR_BIN          = ".bin"
	DIR_INSTALL      = filepath.Join("/usr", "local", "bin", "godibs")
	PROFILES_BASE    = []string{
		PROJECT_NAME,
		PROJECT_NAME + "-dev",
	}
	ARCHITECTURES = []string{
		"amd64",
		"arm64",
	}
)

func Build() error {
	return utils.Build(COMMAND_GO)
}

func BinaryBuild() error {
	platform := os.Getenv("PLATFORM")
	architecture := os.Getenv("ARCHITECTURE")

	return utils.BinaryBuild(platform, architecture, DIR_BIN, COMMAND_GO, PROJECT_NAME, MAIN_FILE)
}

func BinaryInstall() error {
	platform := os.Getenv("PLATFORM")
	architecture := os.Getenv("ARCHITECTURE")

	return utils.BinaryInstall(platform, architecture, DIR_BIN, PROJECT_NAME, DIR_INSTALL)
}

func Clean() error {
	return utils.Clean(DIR_BIN, PROJECT_NAME)
}

func Start() error {
	return utils.Start(COMMAND_GO, MAIN_FILE)
}

func UnitTests() error {
	return utils.UnitTests(COMMAND_GO)
}

func IntegrationTests() error {
	return utils.IntegrationTests(COMMAND_GO, PROJECT_NAME)
}

func BinaryIntegrationTests() error {
	mg.SerialDeps(BinaryInstall)

	return utils.BinaryIntegrationTests(PROJECT_NAME, DIR_INSTALL)
}

func DockerMultiarchSetup() error {
	return utils.DockerMultiarchSetup()
}

func SkaffoldBuild() error {
	mg.SerialDeps(DockerMultiarchSetup)

	return utils.SkaffoldBuild(ARCHITECTURES, PROFILES_BASE)
}

func DockerManifestBuild() error {
	return utils.DockerManifestBuild(DOCKER_NAMESPACE, PROFILES_BASE, ARCHITECTURES)
}
