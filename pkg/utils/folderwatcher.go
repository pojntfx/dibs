package utils

import (
	fswatch "github.com/andreaskoch/go-fswatch"
	"strings"
)

// FolderWatcher watches a folder for changes and sends a message on every change
type FolderWatcher struct {
	FolderWatcher *fswatch.FolderWatcher // Base FolderWatcher
	WatchGlob     string                 // Glob to watch for changes
	IgnoreDir     string                 // Directory to ignore when watching
}

// Start starts the folder watcher
func (folderWatcher *FolderWatcher) Start() {
	folderWatcher.FolderWatcher = fswatch.NewFolderWatcher(folderWatcher.WatchGlob, true, func(path string) bool { return strings.Contains(path, folderWatcher.IgnoreDir) }, 1)

	folderWatcher.FolderWatcher.Start()
}
