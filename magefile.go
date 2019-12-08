//+build mage

package main

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"os"
)

var (
	ARCHITECTURE = os.Getenv("ARCHITECTURE")

	buildConfigAMD64 = utils.BuildConfigV2{
		Platform: "amd64",

		BuildCommand:       "go build -o .bin/godibs-amd64 main.go",
		BuildDockerfile:    "Dockerfile",
		BuildDockerContext: ".",

		TestUnitCommand:       "go test ./...",
		TestUnitDockerfile:    "Dockerfile.testUnit",
		TestUnitDockerContext: ".",
	}
)

func Build() error {
	return buildConfigAMD64.Build()
}

func BuildInDocker() error {
	return buildConfigAMD64.BuildInDocker()
}

func TestUnit() error {
	return buildConfigAMD64.TestUnit()
}

func TestUnitInDocker() error {
	return buildConfigAMD64.TestUnitInDocker()
}
