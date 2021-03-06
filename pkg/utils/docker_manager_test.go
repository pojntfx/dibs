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
	testTag          = "pojntfx/test-app:linux-amd64"
	testManifestTag  = "pojntfx/test-app"
	testExecLine     = "ls"
	testAssetInImage = "/usr/local/bin/test-app" // Don't `filepath.Join` as this is hard-coded in `dibs.yaml` anyways
	testAssetOut     = filepath.Join(os.TempDir(), "test-app")
)

func isDockerEnabled() bool {
	return os.Getenv("DIBS_DISABLE_DOCKER_DEPENDEND_TESTS") != "1"
}

func enableBuildx() error {
	envVariablesToSet := [][]string{
		{"TARGETPLATFORM", "linux/amd64"},
		{"DOCKER_CLI_EXPERIMENTAL", "enabled"},
		{"DOCKER_BUILDKIT", "1"},
	}
	for _, envVariableToSet := range envVariablesToSet {
		if err := os.Setenv(envVariableToSet[0], envVariableToSet[1]); err != nil {
			return err
		}
	}

	return nil
}

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
	if isDockerEnabled() {
		if err := enableBuildx(); err != nil {
			t.Error(err)
		}

		stdoutChan, stderrChan := make(chan string), make(chan string)

		d := NewDockerManager(testContext, stdoutChan, stderrChan)

		hits := 0
		go func() {
			for {
				select {
				case stdout := <-stdoutChan:
					t.Log("test stdout", stdout)
				case stderr := <-stderrChan:
					t.Log("test stderr", stderr)

					if strings.Contains(stderr, "DONE") || strings.Contains(stderr, "naming to") {
						hits++
					}
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
}

// TestPushDockerManager requires the environment variable below to be set; it is disabled by default.
func TestPushDockerManager(t *testing.T) {
	if os.Getenv("DIBS_DOCKER_PUSH_TEST_ENABLED") == "1" {
		if err := enableBuildx(); err != nil {
			t.Error(err)
		}

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
					t.Log("test stderr", stderr)
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
}

func TestRunDockerManager(t *testing.T) {
	if isDockerEnabled() {
		if err := enableBuildx(); err != nil {
			t.Error(err)
		}

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
					t.Log("test stderr", stderr)
				}
			}
		}()

		if err := d.Build(testDockerfile, testContext, testTag); err != nil {
			t.Error(err)
		}

		if err := d.Run(testTag, testExecLine, false); err != nil {
			t.Error(err)
		}

		if hits < 1 {
			t.Error("Docker output did not match expected output")
		}
	}
}

func TestCopyFromImageDockerManager(t *testing.T) {
	if isDockerEnabled() {
		if err := enableBuildx(); err != nil {
			t.Error(err)
		}

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
}

func TestBuildManifestDockerManager(t *testing.T) {
	if isDockerEnabled() {
		if err := enableBuildx(); err != nil {
			t.Error(err)
		}

		stdoutChan, stderrChan := make(chan string), make(chan string)

		d := NewDockerManager(testContext, stdoutChan, stderrChan)

		hits := 0
		go func() {
			for {
				select {
				case stdout := <-stdoutChan:
					t.Log("test stdout", stdout)

					if strings.Contains(stdout, "Created manifest list") {
						hits++
					}
				case stderr := <-stderrChan:
					t.Log("test stderr", stderr)
				}
			}
		}()

		if err := d.Build(testDockerfile, testContext, testTag); err != nil {
			t.Error(err)
		}

		if err := d.BuildManifest(testManifestTag, []string{testTag}); err != nil {
			t.Error(err)
		}

		if hits < 1 {
			t.Error("Docker output did not match expected output")
		}
	}
}

// TestPushManifestDockerManager requires the environment variable below to be set; it is disabled by default.
func TestPushManifestDockerManager(t *testing.T) {
	if os.Getenv("DIBS_DOCKER_MANIFEST_PUSH_TEST_ENABLED") == "1" {
		if err := enableBuildx(); err != nil {
			t.Error(err)
		}

		stdoutChan, stderrChan := make(chan string), make(chan string)

		d := NewDockerManager(testContext, stdoutChan, stderrChan)

		hits := 0
		go func() {
			for {
				select {
				case stdout := <-stdoutChan:
					t.Log("test stdout", stdout)

					if strings.Contains(stdout, "sha256:") { // Only the push command logs to stdout, so this works
						hits++
					}
				case stderr := <-stderrChan:
					t.Log("test stderr", stderr)
				}
			}
		}()

		if err := d.Build(testDockerfile, testContext, testTag); err != nil {
			t.Error(err)
		}

		if err := d.BuildManifest(testManifestTag, []string{testTag}); err != nil {
			t.Error(err)
		}

		if err := d.PushManifest(testManifestTag); err != nil {
			t.Error(err)
		}

		if hits < 1 {
			t.Error("Docker output did not match expected output")
		}
	}
}
