//+build mage

package main

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"os"
)

var (
	ARCHITECTURE = os.Getenv("ARCHITECTURE")

	buildConfigAMD64 = utils.BuildConfigV2{
		Tag:      "pojntfx/godibs:amd64",
		Platform: "amd64",

		BuildCommand:       "go build -o .bin/godibs-amd64 main.go",
		BuildDockerContext: ".",
		BuildDockerfile:    "Dockerfile",

		BuildDockerTag:           "pojntfx/godibs-builddockerindocker:amd64",
		BuildDockerCommand:       "mage buildInDocker",
		BuildDockerDockerContext: ".",
		BuildDockerDockerfile:    "Dockerfile.docker",

		TestUnitCommand:       "go test ./...",
		TestUnitDockerContext: ".",
		TestUnitDockerfile:    "Dockerfile.testUnit",

		TestIntegrationGoCommand:       "go run main.go --help",
		TestIntegrationGoDockerContext: ".",
		TestIntegrationGoDockerfile:    "Dockerfile.testIntegrationGo",

		TestIntegrationBinaryCommand:       ".bin/godibs-amd64 --help",
		TestIntegrationBinaryDockerContext: ".",
		TestIntegrationBinaryDockerfile:    "Dockerfile.testIntegrationBinary",

		TestIntegrationDockerTag:           "pojntfx/godibs-testintegrationdockerindocker:amd64",
		TestIntegrationDockerCommand:       "docker run pojntfx/godibs:amd64 /usr/local/bin/godibs --help",
		TestIntegrationDockerDockerContext: ".",
		TestIntegrationDockerDockerfile:    "Dockerfile.testIntegrationDocker",
	}
)

func Build() error {
	return buildConfigAMD64.Build()
}

func BuildInDocker() error {
	return buildConfigAMD64.BuildInDocker()
}

func BuildDocker() error {
	return buildConfigAMD64.BuildDocker()
}

func BuildDockerInDocker() error {
	return buildConfigAMD64.BuildDockerInDocker()
}

func TestUnit() error {
	return buildConfigAMD64.TestUnit()
}

func TestUnitInDocker() error {
	return buildConfigAMD64.TestUnitInDocker()
}

func TestIntegrationGo() error {
	return buildConfigAMD64.TestIntegrationGo()
}

func TestIntegrationGoInDocker() error {
	return buildConfigAMD64.TestIntegrationGoInDocker()
}

func TestIntegrationBinary() error {
	return buildConfigAMD64.TestIntegrationBinary()
}

func TestIntegrationBinaryInDocker() error {
	return buildConfigAMD64.TestIntegrationBinaryInDocker()
}

func TestIntegrationDocker() error {
	return buildConfigAMD64.TestIntegrationDocker()
}

func TestIntegrationDockerInDocker() error {
	return buildConfigAMD64.TestIntegrationDockerInDocker()
}
