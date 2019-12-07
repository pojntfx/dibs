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

func Build(goCmd string) error {
	return sh.RunV(goCmd, "build", "./...")
}

func BinaryBuild(platform, architecture, dirBin, goCmd, projectName, mainFile string) error {
	if err := os.MkdirAll(dirBin, 0755); err != nil {
		return err
	}

	return sh.RunWith(map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        platform,
		"GOARCH":      architecture,
	}, goCmd, "build", "-o", filepath.Join(dirBin, projectName+"-"+platform+"-"+architecture), mainFile)
}

func BinaryInstall(platform, architecture, dirBin, projectName, dirInstall string) error {
	from, err := os.Open(filepath.Join(dirBin, projectName+"-"+platform+"-"+architecture))
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dirInstall, os.O_RDWR|os.O_CREATE, 755)
	if err != nil {
		return err
	}
	defer to.Close()

	fmt.Println(filepath.Join(dirBin, projectName+"-"+platform+"-"+architecture))
	if _, err := io.Copy(to, from); err != nil {
		return err
	}

	return nil
}

func Clean(dirBin, projectName string) error {
	binariesToRemove, _ := filepath.Glob(filepath.Join(dirBin, projectName+"-*-*"))
	for _, fileToRemove := range binariesToRemove {
		if err := os.Remove(fileToRemove); err != nil {
			return err
		}
	}

	return nil
}

func Start(goCmd, mainFile string) error {
	return sh.RunV(goCmd, append([]string{"run", mainFile}, os.Args[2:]...)...)
}

func UnitTests(goCmd string) error {
	if err := sh.RunV(goCmd, "test", "--tags", "unit", "./..."); err != nil {
		return err
	}

	log.Info("Passed")

	return nil
}

func IntegrationTests(goCmd, projectName string) error {
	if err := sh.RunV(goCmd, "install", "./..."); err != nil {
		return err
	}

	if err := sh.RunV(projectName, "--help"); err != nil {
		return err
	}

	installPathNonBinary, err := exec.LookPath(projectName)
	if err != nil {
		return err
	}
	os.Remove(installPathNonBinary)

	log.Info("Passed")
	return nil
}

func BinaryIntegrationTests(projectName, dirInstall string) error {
	if err := sh.RunV(projectName, "--help"); err != nil {
		return err
	}

	os.Remove(dirInstall)

	log.Info("Passed")
	return nil

}

func DockerMultiarchSetup() error {
	return sh.RunV("docker", "run", "--rm", "--privileged", "multiarch/qemu-user-static", "--reset", "-p", "yes")
}

func SkaffoldBuild(architectures, profilesBase []string) error {
	var profiles []string

	for _, architecture := range architectures {
		profiles = append(profiles, profilesBase[0]+"-"+architecture)
	}

	sh.RunV("skaffold", "config", "unset", "--global", "default-repo")

	for _, profile := range profiles {
		sh.RunV("skaffold", "build", "-p", profile)
	}

	return nil

}

func DockerManifestBuild(dockerNamespace string, profilesBase, architectures []string) error {
	var cmds []string

	manifestName := dockerNamespace + "/" + profilesBase[0] + ":latest"

	cmds = append(cmds, "manifest", "create", "--amend", manifestName)

	for _, architecture := range architectures {
		cmds = append(cmds, dockerNamespace+"/"+profilesBase[0]+":latest-"+architecture)
	}

	err := sh.RunWith(map[string]string{
		"DOCKER_CLI_EXPERIMENTAL": "enabled",
	}, "docker", cmds...)
	if err != nil {
		return err
	}

	return sh.RunV("docker", "manifest", "push", manifestName)
}
