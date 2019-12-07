//+build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	PROJECT_NAME     = "godibs"
	MAIN_FILE        = "main.go"
	DOCKER_NAMESPACE = "pojntfx"
	COMMAND_GO       = mg.GoCmd()
	DIR_TEMP         = os.TempDir()
	DIR_BIN          = ".bin"
	DIR_INSTALL      = filepath.Join("/usr", "local", "bin", "godibs")
	PROFILES_BASE    = []string{
		PROJECT_NAME,
		PROJECT_NAME + "-dev",
	}
	ARCHITECTURES = []string{
		"amd64",
		"arm64",
	}
)

func Build() error {
	return sh.RunV(COMMAND_GO, "build", "./...")
}

func BinaryBuild() error {
	platform := os.Getenv("PLATFORM")
	architecture := os.Getenv("ARCHITECTURE")

	_, err := os.Stat(DIR_BIN)
	if os.IsExist(err) {
		os.Mkdir(DIR_BIN, 0755)
	}

	return sh.RunWith(map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        platform,
		"GOARCH":      architecture,
	}, COMMAND_GO, "build", "-o", filepath.Join(DIR_BIN, PROJECT_NAME+"-"+platform+"-"+architecture), MAIN_FILE)
}

func BinaryInstall() error {
	platform := os.Getenv("PLATFORM")
	architecture := os.Getenv("ARCHITECTURE")

	from, _ := os.Open(filepath.Join(DIR_BIN, PROJECT_NAME+"-"+platform+"-"+architecture))
	defer from.Close()

	to, _ := os.OpenFile(DIR_INSTALL, os.O_RDWR|os.O_CREATE, 755)
	defer to.Close()

	fmt.Println(filepath.Join(DIR_BIN, PROJECT_NAME+"-"+platform+"-"+architecture))
	_, err := io.Copy(to, from)

	return err
}

func Clean() error {
	binariesToRemove, _ := filepath.Glob(filepath.Join(DIR_BIN, PROJECT_NAME+"-*-*"))
	for _, fileToRemove := range binariesToRemove {
		if err := os.Remove(fileToRemove); err != nil {
			return err
		}
	}

	return nil
}

func Start() error {
	return sh.RunV(COMMAND_GO, append([]string{"run", MAIN_FILE}, os.Args[2:]...)...)
}

func UnitTests() error {
	err := sh.RunV(COMMAND_GO, "test", "--tags", "unit", "./...")
	if err != nil {
		return err
	}
	log.Info("Passed")
	return nil
}

func IntegrationTests() error {
	err := sh.RunV(COMMAND_GO, "install", "./...")
	if err != nil {
		return err
	}

	err = sh.RunV(PROJECT_NAME, "--help")
	if err != nil {
		return err
	}

	installPathNonBinary, err := exec.LookPath(PROJECT_NAME)
	if err != nil {
		return err
	}
	os.Remove(installPathNonBinary)

	log.Info("Passed")
	return nil
}

func BinaryIntegrationTests() error {
	mg.SerialDeps(BinaryInstall)

	err := sh.RunV(PROJECT_NAME, "--help")
	if err != nil {
		return err
	}

	os.Remove(DIR_INSTALL)

	log.Info("Passed")
	return nil
}

func DockerMultiarchSetup() error {
	return sh.RunV("docker", "run", "--rm", "--privileged", "multiarch/qemu-user-static", "--reset", "-p", "yes")
}

func SkaffoldBuild() error {
	mg.SerialDeps(DockerMultiarchSetup)
	var profiles []string

	for _, architecture := range ARCHITECTURES {
		profiles = append(profiles, PROFILES_BASE[0]+"-"+architecture)
	}

	sh.RunV("skaffold", "config", "unset", "--global", "default-repo")

	for _, profile := range profiles {
		sh.RunV("skaffold", "build", "-p", profile)
	}

	return nil
}

func DockerManifestBuild() error {
	var cmds []string

	manifestName := DOCKER_NAMESPACE + "/" + PROFILES_BASE[0] + ":latest"

	cmds = append(cmds, "manifest", "create", "--amend", manifestName)

	for _, architecture := range ARCHITECTURES {
		cmds = append(cmds, DOCKER_NAMESPACE+"/"+PROFILES_BASE[0]+":latest-"+architecture)
	}

	err := sh.RunWith(map[string]string{
		"DOCKER_CLI_EXPERIMENTAL": "enabled",
	}, "docker", cmds...)
	if err != nil {
		return err
	}

	return sh.RunV("docker", "manifest", "push", manifestName)
}
