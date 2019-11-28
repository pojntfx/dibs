package main

import (
	redis "github.com/go-redis/redis/v7"
	"github.com/pojntfx/godibs/src/lib/common"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	REDIS_URL            = os.Getenv("REDIS_URL")
	REDIS_CHANNEL_PREFIX = os.Getenv("REDIS_CHANNEL_PREFIX")
	GIT_DIR              = os.Getenv("GIT_DIR")
	GIT_NAME             = os.Getenv("GIT_NAME")
	GIT_EMAIL            = os.Getenv("GIT_EMAIL")
)

func main() {
	r := common.GetNewRedisClient(REDIS_URL)

	StartDirectoryCreationWorker(r, REDIS_CHANNEL_PREFIX, common.REDIS_CHANNEL_MODULE_REGISTERED, GIT_DIR)
}

func ParseModuleFromMessage(m string) (name, timestamp string) {
	res := strings.Split(m, "@")
	return res[0], res[1]
}

func GetPathForModule(baseDir, m string) string {
	return filepath.Join(append([]string{baseDir, "repositories"}, strings.Split(m, "/")...)...)
}

func StartDirectoryCreationWorker(r *redis.Client, prefix, channel, baseDir string) error {
	p := r.Subscribe(prefix + ":" + channel)
	defer p.Close()

	_, err := p.Receive()
	if err != nil {
		return err
	}

	c := p.Channel()

	var wg sync.WaitGroup

	log.Info("Starting directory creation worker ...")
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for m := range c {
			n, t := ParseModuleFromMessage(m.Payload)
			log.Info("Creating directory", rz.String("moduleName", n), rz.String("modulePushedTimestamp", t))

			path := GetPathForModule(baseDir, n)

			err = os.RemoveAll(path)
			if err != nil {
				panic(err)
			}

			err = os.MkdirAll(path, 0777)
			if err != nil {
				panic(err)
			}

		}

		wg.Done()
	}(&wg)

	wg.Wait()

	return nil
}
