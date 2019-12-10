//+build mage

package main

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"os"
)

var (
	PLATFORM = os.Getenv("TARGETPLATFORM")

	buildConfigs = utils.BuildConfigCollection{
		ManifestTag: "pojntfx/godibs:latest",

		BuildConfigs: []utils.BuildConfig{
			utils.BuildConfig{
				Tag:      "pojntfx/godibs:linux-amd64",
				Platform: "linux/amd64",

				BinaryInContainerPath: "/usr/local/bin/godibs",
				BinaryDistPath:        ".bin/godibs-linux-amd64",
				CleanGlob:             ".bin/godibs-linux-amd64",

				BuildCommand:       "go build -o .bin/godibs-linux-amd64 main.go",
				BuildDockerContext: ".",
				BuildDockerfile:    "Dockerfile",

				BuildDockerTag:           "pojntfx/godibs-builddockerindocker:linux-amd64",
				BuildDockerCommand:       "mage buildInDocker",
				BuildDockerDockerContext: ".",
				BuildDockerDockerfile:    "Dockerfile.docker",

				TestUnitCommand:       "go test ./...",
				TestUnitDockerContext: ".",
				TestUnitDockerfile:    "Dockerfile.testUnit",

				TestIntegrationGoCommand:       "go run main.go --help",
				TestIntegrationGoDockerContext: ".",
				TestIntegrationGoDockerfile:    "Dockerfile.testIntegrationGo",

				TestIntegrationBinaryCommand:       ".bin/godibs-linux-amd64 --help",
				TestIntegrationBinaryDockerContext: ".",
				TestIntegrationBinaryDockerfile:    "Dockerfile.testIntegrationBinary",

				TestIntegrationDockerTag:           "pojntfx/godibs-testintegrationdockerindocker:linux-amd64",
				TestIntegrationDockerCommand:       "docker run --platform linux/amd64 pojntfx/godibs:linux-amd64 /usr/local/bin/godibs --help",
				TestIntegrationDockerDockerContext: ".",
				TestIntegrationDockerDockerfile:    "Dockerfile.testIntegrationDocker",
			},
			utils.BuildConfig{
				Tag:      "pojntfx/godibs:linux-arm64",
				Platform: "linux/arm64",

				BinaryInContainerPath: "/usr/local/bin/godibs",
				BinaryDistPath:        ".bin/godibs-linux-arm64",
				CleanGlob:             ".bin/godibs-linux-arm64",

				BuildCommand:       "go build -o .bin/godibs-linux-arm64 main.go",
				BuildDockerContext: ".",
				BuildDockerfile:    "Dockerfile",

				BuildDockerTag:           "pojntfx/godibs-builddockerindocker:linux-arm64",
				BuildDockerCommand:       "mage buildInDocker",
				BuildDockerDockerContext: ".",
				BuildDockerDockerfile:    "Dockerfile.docker",

				TestUnitCommand:       "go test ./...",
				TestUnitDockerContext: ".",
				TestUnitDockerfile:    "Dockerfile.testUnit",

				TestIntegrationGoCommand:       "go run main.go --help",
				TestIntegrationGoDockerContext: ".",
				TestIntegrationGoDockerfile:    "Dockerfile.testIntegrationGo",

				TestIntegrationBinaryCommand:       ".bin/godibs-linux-arm64 --help",
				TestIntegrationBinaryDockerContext: ".",
				TestIntegrationBinaryDockerfile:    "Dockerfile.testIntegrationBinary",

				TestIntegrationDockerTag:           "pojntfx/godibs-testintegrationdockerindocker:linux-arm64",
				TestIntegrationDockerCommand:       "docker run --platform linux/arm64 pojntfx/godibs:linux-arm64 /usr/local/bin/godibs --help",
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

func PushDockerImage() error {
	return buildConfigs.PushDockerImage(PLATFORM)
}

func GetBinaryFromDockerImage() error {
	return buildConfigs.GetBinaryFromDockerImage(PLATFORM)
}

func Clean() error {
	return buildConfigs.Clean(PLATFORM)
}

func BuildDockerManifest() error {
	return buildConfigs.BuildDockerManifest()
}

func PushDockerManifest() error {
	return buildConfigs.PushDockerManifest()
}

func BuildAll() error {
	return buildConfigs.BuildAll()
}

func BuildInDockerAll() error {
	return buildConfigs.BuildInDockerAll()
}

func BuildDockerAll() error {
	return buildConfigs.BuildDockerAll()
}

func BuildDockerInDockerAll() error {
	return buildConfigs.BuildDockerInDockerAll()
}

func TestUnitAll() error {
	return buildConfigs.TestUnitAll()
}

func TestUnitInDockerAll() error {
	return buildConfigs.TestUnitInDockerAll()
}

func TestIntegrationGoAll() error {
	return buildConfigs.TestIntegrationGoAll()
}

func TestIntegrationGoInDockerAll() error {
	return buildConfigs.TestIntegrationGoInDockerAll()
}

func TestIntegrationBinaryAll() error {
	return buildConfigs.TestIntegrationBinaryAll()
}

func TestIntegrationBinaryInDockerAll() error {
	return buildConfigs.TestIntegrationBinaryInDockerAll()
}

func TestIntegrationDockerAll() error {
	return buildConfigs.TestIntegrationDockerAll()
}

func TestIntegrationDockerInDockerAll() error {
	return buildConfigs.TestIntegrationDockerInDockerAll()
}

func PushDockerImageAll() error {
	return buildConfigs.PushDockerImageAll()
}

func GetBinaryFromDockerImageAll() error {
	return buildConfigs.GetBinaryFromDockerImageAll()
}

func CleanAll() error {
	return buildConfigs.CleanAll()
}
