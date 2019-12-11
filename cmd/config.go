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
				BuildCleanGlob:     ".bin/godibs-linux-amd64",
				BuildDockerContext: ".",
				BuildDockerfile:    "Dockerfile",

				BuildImageCommand: "mage buildInDocker",

				TestUnitCommand:       "go test ./...",
				TestUnitDockerContext: ".",
				TestUnitDockerfile:    "Dockerfile.testUnit",
				TestUnitImageTag:      "pojntfx/godibs-unittest-go:linux-amd64",

				TestIntegrationLangCommand:       "go run main.go --help",
				TestIntegrationLangDockerContext: ".",
				TestIntegrationLangDockerfile:    "Dockerfile.testIntegrationLang",
				TestIntegrationLangImageTag:      "pojntfx/godibs-integrationtest-lang:linux-amd64",

				TestIntegrationBinaryCommand:       ".bin/godibs-linux-amd64 --help",
				TestIntegrationBinaryDockerContext: ".",
				TestIntegrationBinaryDockerfile:    "Dockerfile.testIntegrationBinary",
				TestIntegrationBinaryImageTag:      "pojntfx/godibs-integrationtest-binary:linux-amd64",

				TestIntegrationImageCommand:       "docker run --platform linux/amd64 pojntfx/godibs:linux-amd64 /usr/local/bin/godibs --help",
				TestIntegrationImageDockerContext: ".",
				TestIntegrationImageDockerfile:    "Dockerfile.testIntegrationImage",
				TestIntegrationImageImageTag:      "pojntfx/godibs-integrationtest-image:linux-amd64",
			},
			utils.BuildConfig{
				Tag:      "pojntfx/godibs:linux-arm64",
				Platform: "linux/arm64",

				BinaryInContainerPath: "/usr/local/bin/godibs",
				BinaryDistPath:        ".bin/godibs-linux-arm64",
				CleanGlob:             ".bin/godibs-linux-arm64",

				BuildCommand:       "go build -o .bin/godibs-linux-arm64 main.go",
				BuildCleanGlob:     ".bin/godibs-linux-arm64",
				BuildDockerContext: ".",
				BuildDockerfile:    "Dockerfile",

				BuildImageCommand: "mage buildInDocker",

				TestUnitCommand:       "go test ./...",
				TestUnitDockerContext: ".",
				TestUnitDockerfile:    "Dockerfile.testUnit",
				TestUnitImageTag:      "pojntfx/godibs-unittest-go:linux-arm64",

				TestIntegrationLangCommand:       "go run main.go --help",
				TestIntegrationLangDockerContext: ".",
				TestIntegrationLangDockerfile:    "Dockerfile.testIntegrationLang",
				TestIntegrationLangImageTag:      "pojntfx/godibs-integrationtest-lang:linux-arm64",

				TestIntegrationBinaryCommand:       ".bin/godibs-linux-arm64 --help",
				TestIntegrationBinaryDockerContext: ".",
				TestIntegrationBinaryDockerfile:    "Dockerfile.testIntegrationBinary",
				TestIntegrationBinaryImageTag:      "pojntfx/godibs-integrationtest-binary:linux-arm64",

				TestIntegrationImageCommand:       "docker run --platform linux/arm64 pojntfx/godibs:linux-arm64 /usr/local/bin/godibs --help",
				TestIntegrationImageDockerContext: ".",
				TestIntegrationImageDockerfile:    "Dockerfile.testIntegrationImage",
				TestIntegrationImageImageTag:      "pojntfx/godibs-integrationtest-image:linux-arm64",
			}}}
)
