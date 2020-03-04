package utils

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	testRepoRoot = filepath.Join(testContext, "..", "..")
)

func TestCreateBinaryManager(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	b := NewBinaryManager(testContext, stdoutChan, stderrChan)

	if b == nil {
		t.Error("New binary manager is nil")
	}

	if b.dir != testContext {
		t.Error("dir not set correctly")
	}

	if b.stdoutChan != stdoutChan {
		t.Error("stdoutChan not set correctly")
	}

	if b.stderrChan != stderrChan {
		t.Error("stderrChan not correctly")
	}
}

// TestPushBinaryManager requires the environment variables below to be set; it is disabled by default.
func TestPushBinaryManager(t *testing.T) {
	if os.Getenv("DIBS_BINARY_PUSH_TEST_ENABLED") == "1" {
		if err := enableBuildx(); err != nil {
			t.Error(err)
		}

		stdoutChan, stderrChan := make(chan string), make(chan string)

		b := NewBinaryManager(testContext, stdoutChan, stderrChan)
		d := NewDockerManager(testContext, stdoutChan, stderrChan)

		go func() {
			for {
				select {
				case stdout := <-stdoutChan:
					t.Log("test stdout", stdout)
				case stderr := <-stderrChan:
					t.Error("error while building or pushing binary chart", stderr)
				}
			}
		}()

		if err := d.Build(testDockerfile, testContext, testTag); err != nil {
			t.Error(err)
		}

		if err := b.Push(
			os.Getenv("DIBS_GITHUB_USER_NAME"),
			os.Getenv("DIBS_GITHUB_TOKEN"),
			os.Getenv("DIBS_GITHUB_REPOSITORY"),
			testRepoRoot,
			testAssetOut,
		); err != nil {
			t.Error(err)
		}
	}
}
