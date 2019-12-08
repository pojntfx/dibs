//+build mage

package main

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"os"
)

var (
	PLATFORM = os.Getenv("TARGETPLATFORM")

	buildConfigs = utils.BuildConfigCollectionV2{
		BuildConfigs: []utils.BuildConfigV2{
			utils.BuildConfigV2{
				Tag:      "pojntfx/godibs:amd64",
				Platform: "linux/amd64",

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
				TestIntegrationDockerCommand:       "docker run --platform amd64 pojntfx/godibs:amd64 /usr/local/bin/godibs --help",
				TestIntegrationDockerDockerContext: ".",
				TestIntegrationDockerDockerfile:    "Dockerfile.testIntegrationDocker",
			}}}
)

func Build() error {
	return buildConfigs.Build(PLATFORM)
}

func BuildInDocker() error {
	return buildConfigs.BuildInDocker(PLATFORM)
}

func BuildDocker() error {
	return buildConfigs.BuildDocker(PLATFORM)
}

func BuildDockerInDocker() error {
	return buildConfigs.BuildDockerInDocker(PLATFORM)
}

func TestUnit() error {
	return buildConfigs.TestUnit(PLATFORM)
}

func TestUnitInDocker() error {
	return buildConfigs.TestUnitInDocker(PLATFORM)
}

func TestIntegrationGo() error {
	return buildConfigs.TestIntegrationGo(PLATFORM)
}

func TestIntegrationGoInDocker() error {
	return buildConfigs.TestIntegrationGoInDocker(PLATFORM)
}

func TestIntegrationBinary() error {
	return buildConfigs.TestIntegrationBinary(PLATFORM)
}

func TestIntegrationBinaryInDocker() error {
	return buildConfigs.TestIntegrationBinaryInDocker(PLATFORM)
}

func TestIntegrationDocker() error {
	return buildConfigs.TestIntegrationDocker(PLATFORM)
}

func TestIntegrationDockerInDocker() error {
	return buildConfigs.TestIntegrationDockerInDocker(PLATFORM)
}
