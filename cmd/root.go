package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var (
	REDIS_URL    string
	REDIS_PREFIX string
	ON           string
	PLATFORM     string
)

const (
	REDIS_SUFFIX_UP_BUILT        = "up_built"
	REDIS_SUFFIX_UP_TESTED       = "up_tested"
	REDIS_SUFFIX_UP_STARTED      = "up_started"
	REDIS_SUFFIX_UP_REGISTERED   = "up_registered"
	REDIS_SUFFIX_UP_UNREGISTERED = "up_unregistered"
	REDIS_SUFFIX_UP_PUSHED       = "up_pushed"
	REDIS_SUFFIX_DOWN_DOWNLOADED = "down_downloaded"

	ON_NATIVE = "native"
	ON_DOCKER = "docker"

	PLATFORM_ALL     = "all"
	PLATFORM_DEFAULT = PLATFORM_ALL
)

// rootCmd ist the main entry command
var rootCmd = &cobra.Command{
	Use:   "godibs",
	Short: "System for distributed multi-module, multi-architecture development with Go",
}

// init maps the flags to the config
func init() {
	rootCmd.PersistentFlags().StringVar(&REDIS_URL, "redis-url", "localhost:6379", "URL of the Redis instance to use")
	rootCmd.PersistentFlags().StringVar(&REDIS_PREFIX, "redis-prefix", "godibs", "Redis channel prefix")
	rootCmd.PersistentFlags().StringVar(&ON, "on", ON_NATIVE, "System to run on (native|docker)")
	rootCmd.PersistentFlags().StringVar(&PLATFORM, "platform", PLATFORM_DEFAULT, "Platform specified in configuration to use (\""+PLATFORM_ALL+"\" to run for every platform)")
}

// Execute starts the main entry command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
	}

	if !(ON == ON_NATIVE || ON == ON_DOCKER) {
		log.Fatal("Unsupported value for --on, must be either \""+ON_NATIVE+"\" or \""+ON_DOCKER+"\"", rz.String("--on", ON))
	}
}
