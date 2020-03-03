package utils

import "testing"

func TestCreateDockerManager(t *testing.T) {
	d := NewDockerManager()

	if d == nil {
		t.Error("New Docker manager is nil")
	}
}
