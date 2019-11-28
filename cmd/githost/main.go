package main

import (
	"github.com/pojntfx/godibs/src/lib/common"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"os"
	"strings"
)

var (
	REDIS_URL            = os.Getenv("REDIS_URL")
	REDIS_CHANNEL_PREFIX = os.Getenv("REDIS_CHANNEL_PREFIX")
)

func main() {
	r := common.GetNewRedisClient(REDIS_URL)

	p := r.Subscribe(REDIS_CHANNEL_PREFIX + ":" + common.REDIS_CHANNEL_MODULE_REGISTERED)
	defer p.Close()

	_, err := p.Receive()
	if err != nil {
		panic(err)
	}

	c := p.Channel()

	log.Info("Starting main loop ...")
	for m := range c {
		n, t := parseModuleFromMessage(m.Payload)
		log.Info("Module registered", rz.String("moduleName", n), rz.String("modulePushedTimestamp", t))
	}
}

func parseModuleFromMessage(m string) (name, timestamp string) {
	res := strings.Split(m, "@")
	return res[0], res[1]
}
