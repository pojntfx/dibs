package utils

import (
	"github.com/radovskyb/watcher"
	"regexp"
	"time"
)

// FolderWatcher watches a folder for changes and sends sends an event on every change
type FolderWatcher struct {
	FolderWatcher *watcher.Watcher // Base FolderWatcher
	WatchDir      string           // Directory to watch for changes
	IgnoreRegex   string           // Regex of paths to ignore
}

// Start starts the folder watcher
func (folderWatcher *FolderWatcher) Start(errorHandler func(err error), eventHandler func(event watcher.Event)) error {
	folderWatcher.FolderWatcher = watcher.New()
	defer folderWatcher.FolderWatcher.Close()

	folderWatcher.FolderWatcher.SetMaxEvents(1)

	go func() {
		for {
			select {
			case event := <-folderWatcher.FolderWatcher.Event:
				ignore, err := regexp.MatchString(folderWatcher.IgnoreRegex, event.Path)
				if err != nil {
					errorHandler(err)
				}
				if !ignore {
					eventHandler(event)
				}
			case err := <-folderWatcher.FolderWatcher.Error:
				errorHandler(err)
			case <-folderWatcher.FolderWatcher.Closed:
				return
			}
		}
	}()

	if err := folderWatcher.FolderWatcher.AddRecursive(folderWatcher.WatchDir); err != nil {
		return err
	}

	return folderWatcher.FolderWatcher.Start(time.Millisecond * 100)
}
