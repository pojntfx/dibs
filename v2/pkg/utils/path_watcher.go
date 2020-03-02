package utils

import (
	"github.com/radovskyb/watcher"
	"regexp"
	"time"
)

// PathWatcher watches for changes in a directory
type PathWatcher struct {
	pathWatch, pathInclude string
	eventChan              chan string
}

// NewPathWatcher creates a new PathWatcher
func NewPathWatcher(pathWatch, pathInclude string, eventChan chan string) *PathWatcher {
	return &PathWatcher{
		pathWatch:   pathWatch,
		pathInclude: pathInclude,
		eventChan:   eventChan,
	}
}

// Start starts the PathWatcher
func (p *PathWatcher) Start() error {
	watch := watcher.New()

	pathIncludeRegex := regexp.MustCompile(p.pathInclude)
	watch.AddFilterHook(watcher.RegexFilterHook(pathIncludeRegex, true))

	if err := watch.AddRecursive(p.pathWatch); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-watch.Event:
				p.eventChan <- event.Path
			case <-watch.Closed:
				return
			}
		}
	}()

	return watch.Start(time.Millisecond * 100)
}
