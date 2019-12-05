package utils

import (
	fswatch "github.com/andreaskoch/go-fswatch"
	"regexp"
)

// FolderWatcher watches a folder for changes and sends a message on every change
type FolderWatcher struct {
	FolderWatcher *fswatch.FolderWatcher // Base FolderWatcher
	WatchDir      string                 // Directory to watch for changes
	IgnoreRegex   string                 // Regex of paths to ignore
}

// Start starts the folder watcher
func (folderWatcher *FolderWatcher) Start() {
	folderWatcher.FolderWatcher = fswatch.NewFolderWatcher(folderWatcher.WatchDir, true, func(path string) bool {
		matched, _ := regexp.Match(folderWatcher.IgnoreRegex, []byte(path))
		return matched
	}, 1)

	folderWatcher.FolderWatcher.Start()
}
