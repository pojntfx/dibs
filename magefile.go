//+build mage

package main

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"os"
	"path/filepath"
)

var buildConfigAMD64 = utils.BuildConfig{
	DockerContext:                ".",
	Dockerfile:                   "Dockerfile.amd64",
	Architecture:                 "amd64",
	Tag:                          "pojntfx/godibs:amd64",
	BinaryInContainerPath:        "/usr/local/bin/godibs",
	BinaryDistPath:               filepath.Join(".bin", "godibs-amd64"),
	IntegrationTestCommandBinary: ".bin/godibs-amd64 --help",
	IntegrationTestCommandDocker: "docker run pojntfx/godibs:amd64",
}

var buildConfigARM64 = utils.BuildConfig{
	DockerContext:                ".",
	Dockerfile:                   "Dockerfile.arm64",
	Architecture:                 "arm64",
	Tag:                          "pojntfx/godibs:arm64",
	BinaryInContainerPath:        "/usr/local/bin/godibs",
	BinaryDistPath:               filepath.Join(".bin", "godibs-arm64"),
	IntegrationTestCommandBinary: ".bin/godibs-arm64 --help",
	IntegrationTestCommandDocker: "docker run pojntfx/godibs:arm64",
}

var buildConfigARM = utils.BuildConfig{
	DockerContext:                ".",
	Dockerfile:                   "Dockerfile.arm",
	Architecture:                 "arm",
	Tag:                          "pojntfx/godibs:arm",
	BinaryInContainerPath:        "/usr/local/bin/godibs",
	BinaryDistPath:               filepath.Join(".bin", "godibs-arm"),
	IntegrationTestCommandBinary: ".bin/godibs-arm --help",
	IntegrationTestCommandDocker: "docker run pojntfx/godibs:arm",
}

var buildConfigCollection = utils.BuildConfigCollection{
	Tag:                    "pojntfx/godibs",
	UnitTestCommand:        "go test ./...",
	IntegrationTestCommand: "go run main.go server --help",
	CleanGlob:              ".bin",
	BuildConfigs: []utils.BuildConfig{
		buildConfigAMD64,
		buildConfigARM64,
		buildConfigARM,
	},
}

var ARCHITECTURE = os.Getenv("ARCHITECTURE")

func BuildDockerImage() error {
	return buildConfigCollection.BuildDockerImage(ARCHITECTURE)
}

func PushDockerImage() error {
	return buildConfigCollection.PushDockerImage(ARCHITECTURE)
}

func GetBinaryFromDockerContainer() error {
	return buildConfigCollection.GetBinaryFromDockerContainer(ARCHITECTURE)
}

func IntegrationTestDocker() error {
	return buildConfigCollection.IntegrationTestDocker(ARCHITECTURE)
}

func IntegrationTestBinary() error {
	return buildConfigCollection.IntegrationTestBinary(ARCHITECTURE)
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

func BuildDockerImagesAll() error {
	return buildConfigCollection.BuildDockerImagesAll()
}

func PushDockerImagesAll() error {
	return buildConfigCollection.PushDockerImagesAll()
}

func BuildDockerManifest() error {
	return buildConfigCollection.BuildDockerManifest()
}

func PushDockerManifest() error {
	return buildConfigCollection.PushDockerManifest()
}

func GetBinariesFromDockerContainersAll() error {
	return buildConfigCollection.GetBinariesFromDockerContainersAll()
}
