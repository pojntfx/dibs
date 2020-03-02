package utils

import (
	"testing"
)

func TestCreatePathWatcher(t *testing.T) {
	pathWatch := "."
	pathInclude := ".*"
	eventChan := make(chan string)

	w := NewPathWatcher(pathWatch, pathInclude, eventChan)

	if w == nil {
		t.Error("New path watcher is nil")
	}

	if w.pathWatch != pathWatch {
		t.Error("pathWatch not set correctly")
	}

	if w.pathInclude != pathInclude {
		t.Error("pathInclude not set correctly")
	}

	if w.eventChan != eventChan {
		t.Error("eventChan not set correctly")
	}
}

func TestStartPathWatcher(t *testing.T) {
	pathWatch := "."
	pathInclude := ".*"
	eventChan := make(chan string)

	w := NewPathWatcher(pathWatch, pathInclude, eventChan)

	// FIXME: Add tests that test if the watcher actually gets a file change, creation and deletion event etc.
	go func() {
		if err := w.Start(); err != nil {
			t.Error(err)
		}
	}()
}
