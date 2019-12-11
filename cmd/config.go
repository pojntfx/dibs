package cmd

import "github.com/pojntfx/godibs/pkg/utils"

var (
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
