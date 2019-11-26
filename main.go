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
	REDIS_CHANNEL               = os.Getenv("REDIS_CHANNEL")
	SYNC_MODULE_PUSH_MOD_FILE   = os.Getenv("SYNC_MODULE_PUSH_MOD_FILE")
	SYNC_MODULE_PUSH_WATCH_GLOB = os.Getenv("SYNC_MODULE_PUSH_WATCH_GLOB")
)

func main() {
	w := fswatch.NewFolderWatcher(SYNC_MODULE_PUSH_WATCH_GLOB, true, func(path string) bool { return false }, 1)
	w.Start()

	r := redis.NewClient(&redis.Options{
		Addr: REDIS_URL,
	})

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

	for w.IsRunning() {
		select {
		case <-w.ChangeDetails():
			t := time.Now().UnixNano()
			e := m + "@" + strconv.Itoa(int(t))

			r.Publish(REDIS_CHANNEL, e)

			fmt.Println(m)
		}
	}
}
