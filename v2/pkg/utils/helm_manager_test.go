package utils

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	testHelmChartSrc  = filepath.Join(testContext, "charts", "test-app")
	testHelmChartDist = filepath.Join(os.TempDir(), "test-app-charts")
)

func TestCreateHelmManager(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	h := NewHelmManager(testContext, stdoutChan, stderrChan)

	if h == nil {
		t.Error("New Helm manager is nil")
	}

	if h.dir != testContext {
		t.Error("dir not set correctly")
	}

	if h.stdoutChan != stdoutChan {
		t.Error("stdoutChan not set correctly")
	}

	if h.stderrChan != stderrChan {
		t.Error("stderrChan not correctly")
	}
}

func TestBuildHelmManager(t *testing.T) {
	stdoutChan, stderrChan := make(chan string), make(chan string)

	h := NewHelmManager(testContext, stdoutChan, stderrChan)

	if err := os.RemoveAll(testHelmChartDist); err != nil {
		t.Error(err)
	}

	if err := os.MkdirAll(testHelmChartDist, 0777); err != nil {
		t.Error(err)
	}

	go func() {
		for {
			select {
			case stdout := <-stdoutChan:
				t.Log("test stdout", stdout)
			case stderr := <-stderrChan:
				t.Error("error while building Helm chart", stderr)
			}
		}
	}()

	if err := h.Build(testHelmChartSrc, testHelmChartDist); err != nil {
		t.Error(err)
	}

	matches, err := filepath.Glob(filepath.Join(testHelmChartDist, "*.tgz"))
	if err != nil {
		t.Error(err)
	}

	if len(matches) == 0 {
		t.Error("Built Helm chart could not be found")
	}
}
