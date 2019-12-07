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

func BuildDockerImageAMD64() error {
	return buildConfigAMD64.BuildDockerImage()
}

func BuildDockerImageARM64() error {
	return buildConfigARM64.BuildDockerImage()
}

func BuildDockerImageARM() error {
	return buildConfigARM.BuildDockerImage()
}

func PushDockerImageAMD64() error {
	return buildConfigAMD64.PushDockerImage()
}

func PushDockerImageARM64() error {
	return buildConfigARM64.PushDockerImage()
}

func PushDockerImageARM() error {
	return buildConfigARM.PushDockerImage()
}

func GetBinaryFromDockerContainerAMD64() error {
	return buildConfigAMD64.GetBinaryFromDockerContainer()
}

func GetBinaryFromDockerContainerARM64() error {
	return buildConfigARM64.GetBinaryFromDockerContainer()
}

func GetBinaryFromDockerContainerARM() error {
	return buildConfigARM.GetBinaryFromDockerContainer()
}

func IntegrationTestDockerAMD64() error {
	return buildConfigAMD64.IntegrationTestDocker()
}

func IntegrationTestDockerARM64() error {
	return buildConfigARM64.IntegrationTestDocker()
}

func IntegrationTestDockerARM() error {
	return buildConfigARM.IntegrationTestDocker()
}

func IntegrationTestBinaryAMD64() error {
	return buildConfigAMD64.IntegrationTestDocker()
}

func IntegrationTestBinaryARM64() error {
	return buildConfigARM64.IntegrationTestDocker()
}

func IntegrationTestBinaryARM() error {
	return buildConfigARM.IntegrationTestDocker()
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

func BuildAllDockerImages() error {
	return buildConfigCollection.BuildAllDockerImages()
}

func PushAllDockerImages() error {
	return buildConfigCollection.PushAllDockerImages()
}

func BuildDockerManifest() error {
	return buildConfigCollection.BuildDockerManifest()
}

func PushDockerManifest() error {
	return buildConfigCollection.PushDockerManifest()
}

func GetAllBinariesFromDockerContainers() error {
	return buildConfigCollection.GetAllBinariesFromDockerContainers()
}
