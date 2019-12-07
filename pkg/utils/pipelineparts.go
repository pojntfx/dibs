package utils

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type BuildConfiguration struct {
	ProjectName      string   // Name of the project
	MainFile         string   // Main file of the project
	DockerRepoPrefix string   // Docker repo prefix
	GoCmd            string   // Location to the Go binary
	DirTemp          string   // Temporary directory to work in
	DirBin           string   // Directory to place the built binaries in
	DirInstall       string   // Directory to install the built binaries to
	ProfilesBase     []string // Base Skaffold profiles
	Architectures    []string // Architectures to build Docker images for
}

func (config *BuildConfiguration) Build() error {
	return sh.RunV(config.GoCmd, "build", "./...")
}

func (config *BuildConfiguration) BinaryBuild(platform, architecture string) error {
	if err := os.MkdirAll(config.DirBin, 0755); err != nil {
		return err
	}

	return sh.RunWith(map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        platform,
		"GOARCH":      architecture,
	}, config.GoCmd, "build", "-o", filepath.Join(config.DirBin, config.ProjectName+"-"+platform+"-"+architecture), config.MainFile)
}

func (config *BuildConfiguration) BinaryInstall(platform, architecture string) error {
	from, err := os.Open(filepath.Join(config.DirBin, config.ProjectName+"-"+platform+"-"+architecture))
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(config.DirInstall, os.O_RDWR|os.O_CREATE, 755)
	if err != nil {
		return err
	}
	defer to.Close()

	fmt.Println(filepath.Join(config.DirBin, config.ProjectName+"-"+platform+"-"+architecture))
	if _, err := io.Copy(to, from); err != nil {
		return err
	}

	return nil
}

func (config *BuildConfiguration) Clean() error {
	binariesToRemove, _ := filepath.Glob(filepath.Join(config.DirBin, config.ProjectName+"-*-*"))
	for _, fileToRemove := range binariesToRemove {
		if err := os.Remove(fileToRemove); err != nil {
			return err
		}
	}

	return nil
}

func (config *BuildConfiguration) Start() error {
	return sh.RunV(config.GoCmd, append([]string{"run", config.MainFile}, os.Args[2:]...)...)
}

func (config *BuildConfiguration) UnitTests() error {
	if err := sh.RunV(config.GoCmd, "test", "--tags", "unit", "./..."); err != nil {
		return err
	}

	log.Info("Passed")

	return nil
}

func (config *BuildConfiguration) IntegrationTests() error {
	if err := sh.RunV(config.GoCmd, "install", "./..."); err != nil {
		return err
	}

	if err := sh.RunV(config.ProjectName, "--help"); err != nil {
		return err
	}

	installPathNonBinary, err := exec.LookPath(config.ProjectName)
	if err != nil {
		return err
	}
	os.Remove(installPathNonBinary)

	log.Info("Passed")
	return nil
}

func (config *BuildConfiguration) BinaryIntegrationTests() error {
	if err := sh.RunV(config.ProjectName, "--help"); err != nil {
		return err
	}

	os.Remove(config.DirInstall)

	log.Info("Passed")
	return nil

}

func (config *BuildConfiguration) DockerMultiarchSetup() error {
	return sh.RunV("docker", "run", "--rm", "--privileged", "multiarch/qemu-user-static", "--reset", "-p", "yes")
}

func (config *BuildConfiguration) SkaffoldBuild() error {
	var profiles []string

	for _, architecture := range config.Architectures {
		profiles = append(profiles, config.ProfilesBase[0]+"-"+architecture)
	}

	sh.RunV("skaffold", "config", "unset", "--global", "default-repo")

	for _, profile := range profiles {
		sh.RunV("skaffold", "build", "-p", profile)
	}

	return nil

}

func (config *BuildConfiguration) DockerManifestBuild() error {
	var cmds []string

	manifestName := config.DockerRepoPrefix + "/" + config.ProfilesBase[0] + ":latest"

	cmds = append(cmds, "manifest", "create", "--amend", manifestName)

	for _, architecture := range config.Architectures {
		cmds = append(cmds, config.DockerRepoPrefix+"/"+config.ProfilesBase[0]+":latest-"+architecture)
	}

	err := sh.RunWith(map[string]string{
		"DOCKER_CLI_EXPERIMENTAL": "enabled",
	}, "docker", cmds...)
	if err != nil {
		return err
	}

	return sh.RunV("docker", "manifest", "push", manifestName)
}
