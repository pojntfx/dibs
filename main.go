package main

import (
	"fmt"
	fswatch "github.com/andreaskoch/go-fswatch"
	"github.com/go-redis/redis/v7"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	REDIS_URL                   = os.Getenv("REDIS_URL")
	REDIS_CHANNEL_PREFIX        = os.Getenv("REDIS_CHANNEL_PREFIX")
	SYNC_MODULE_PUSH_MOD_FILE   = os.Getenv("SYNC_MODULE_PUSH_MOD_FILE")
	SYNC_MODULE_PUSH_WATCH_GLOB = os.Getenv("SYNC_MODULE_PUSH_WATCH_GLOB")
)

func main() {
	f, err := ioutil.ReadFile(SYNC_MODULE_PUSH_MOD_FILE)
	if err != nil {
		panic(err)
	}

	var m string

	for _, line := range strings.Split(string(f), "\n") {
		if strings.Contains(line, "module") {
			m = strings.Split(line, "module ")[1]
			break
		}
	}

	r := redis.NewClient(&redis.Options{
		Addr: REDIS_URL,
	})
	defer r.Publish(REDIS_CHANNEL_PREFIX+":"+"module_unregistered", withTimestamp(m))
	r.Publish(REDIS_CHANNEL_PREFIX+":"+"module_registered", withTimestamp(m))

	w := fswatch.NewFolderWatcher(SYNC_MODULE_PUSH_WATCH_GLOB, true, func(path string) bool { return false }, 1)
	w.Start()

	for w.IsRunning() {
		select {
		case <-w.ChangeDetails():
			r.Publish(REDIS_CHANNEL_PREFIX+":"+"module_pushed", withTimestamp(m))
			fmt.Println(m)
		}
	}
}

func withTimestamp(m string) string {
	t := time.Now().UnixNano()
	return m + "@" + strconv.Itoa(int(t))
}
