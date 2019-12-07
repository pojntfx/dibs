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
	PLATFORM = os.Getenv("PLATFORM")
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

var (
	ARCHITECTURE = os.Getenv("ARCHITECTURE")
	TAG          = os.Getenv("TAG")
)

func BuildDocker() error {
	return buildDocker(ARCHITECTURE, TAG)
}

func BuildBinary() error {
	mg.SerialDeps(BuildDocker)

	return buildBinary(TAG, "/usr/local/bin/godibs", filepath.Join(".bin", "godibs-"+ARCHITECTURE))
}

func buildDocker(architecture, tag string) error {
	return sh.RunV("docker", "build", "-f", "Dockerfile."+architecture, "-t", tag, ".")
}

func buildBinary(tag, srcDir, distDir string) error {
	id := uuid.New().String()

	out, err := exec.Command("docker", "ps", "-aqf", "name="+id).Output()
	if err != nil {
		return err
	}
	if string(out) != "\n" {
		if err := sh.RunV("docker", "run", "--name", id, tag, "echo"); err != nil {
			return err
		}

		if err := sh.RunV("docker", "cp", id+":"+srcDir, distDir); err != nil {
			return err
		}

		if err := sh.RunV("docker", "rm", "-f", id); err != nil {
			return err
		}
	} else {
		return buildBinary(tag, srcDir, distDir)
	}

	return nil
}
