package main

import (
	"fmt"
	fswatch "github.com/andreaskoch/go-fswatch"
	"github.com/go-redis/redis/v7"
	"os"
)

var (
	REDIS                       = os.Getenv("REDIS")
	SYNC_MODULE_PUSH_WATCH_GLOB = os.Getenv("SYNC_MODULE_PUSH_WATCH_GLOB")
	CHANNEL                     = "godibs:pushed_module"
)

func main() {
	w := fswatch.NewFolderWatcher(SYNC_MODULE_PUSH_WATCH_GLOB, true, func(path string) bool { return false }, 1)
	w.Start()

	r := redis.NewClient(&redis.Options{
		Addr: REDIS,
	})

	for w.IsRunning() {
		select {
		case <-w.ChangeDetails():
			r.Publish(CHANNEL, "test message")
			fmt.Println("Published to channel after file change")
		}
	}
}
