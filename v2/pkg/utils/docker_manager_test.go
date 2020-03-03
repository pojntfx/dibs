package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	testContext = func() string {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		return filepath.Join(pwd, "..", "..", "test-app")
	}()
	testDockerfile   = filepath.Join(testContext, "Dockerfile")
	testTag          = filepath.Join("pojntfx/test-app")
	testExecLine     = "ls"
	testAssetInImage = "/usr/local/bin/test-app" // Don't `filepath.Join` as this is hard-coded in `dibs.yaml` anyways
	testAssetOut     = filepath.Join(os.TempDir(), "test-app")
)

func TestCreateDockerManager(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	d := NewDockerManager(testContext, stdoutChan, stderrChan)

	if d == nil {
		t.Error("New Docker manager is nil")
	}

	if d.dir != testContext {
		t.Error("dir not set correctly")
	}

	if d.stdoutChan != stdoutChan {
		t.Error("stdoutChan not set correctly")
	}

	if d.stderrChan != stderrChan {
		t.Error("stderrChan not correctly")
	}
}

func TestBuildDockerManager(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	d := NewDockerManager(testContext, stdoutChan, stderrChan)

	hits := 0
	go func() {
		for {
			select {
			case stdout := <-stdoutChan:
				t.Log("test stdout", stdout)

				if strings.Contains(stdout, "Successfully built") || strings.Contains(stdout, "Successfully tagged") {
					hits++
				}
			case stderr := <-stderrChan:
				t.Error("error while building Docker image", stderr)
			}
		}
	}()

	if err := d.Build(testDockerfile, testContext, testTag); err != nil {
		t.Error(err)
	}

	if hits < 2 {
		t.Error("Docker output did not match expected output")
	}
}

func TestPushDockerManager(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	d := NewDockerManager(testContext, stdoutChan, stderrChan)

	hits := 0
	go func() {
		for {
			select {
			case stdout := <-stdoutChan:
				t.Log("test stdout", stdout)

				if strings.Contains(stdout, "The push refers to repository") || strings.Contains(stdout, "digest:") {
					hits++
				}
			case stderr := <-stderrChan:
				t.Error("error while building or pushing Docker image", stderr)
			}
		}
	}()

	if err := d.Build(testDockerfile, testContext, testTag); err != nil {
		t.Error(err)
	}

	if err := d.Push(testTag); err != nil {
		t.Error(err)
	}

	if hits < 2 {
		t.Error("Docker output did not match expected output")
	}
}

func TestRunDockerManager(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	d := NewDockerManager(testContext, stdoutChan, stderrChan)

	hits := 0
	go func() {
		for {
			select {
			case stdout := <-stdoutChan:
				t.Log("test stdout", stdout)

				if strings.Contains(stdout, "usr") {
					hits++
				}
			case stderr := <-stderrChan:
				t.Error("error while building or running Docker image", stderr)
			}
		}
	}()

	if err := d.Build(testDockerfile, testContext, testTag); err != nil {
		t.Error(err)
	}

	if err := d.Run(testTag, testExecLine, false); err != nil {
		t.Error(err)
	}

	if hits < 2 {
		t.Error("Docker output did not match expected output")
	}
}

func TestCopyFromImageDockerManager(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	d := NewDockerManager(testContext, stdoutChan, stderrChan)

	if err := d.Build(testDockerfile, testContext, testTag); err != nil {
		t.Error(err)
	}

	if err := os.RemoveAll(testAssetOut); err != nil {
		t.Error(err)
	}

	if err := d.CopyFromImage(testTag, testAssetInImage, testAssetOut); err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(testAssetOut); err != nil {
		t.Error(err)
	}
}