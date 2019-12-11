package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var (
	RedisUrl    string
	RedisPrefix string
	ON          string
	PLATFORM    string
)

const (
	RedisSuffixUpBuilt        = "up_built"
	RedisSuffixUpTested       = "up_tested"
	RedisSuffixUpStarted      = "up_started"
	RedisSuffixUpRegistered   = "up_registered"
	RedisSuffixUpUnregistered = "up_unregistered"
	RedisSuffixUpPushed       = "up_pushed"

	OnNative = "native"
	OnDocker = "docker"

	PlatformAll     = "all"
	PlatformDefault = PlatformAll
)

// rootCmd ist the main entry command
var rootCmd = &cobra.Command{
	Use:   "godibs",
	Short: "System for distributed multi-module, multi-architecture development with Go",
}

// init maps the flags to the config
func init() {
	rootCmd.PersistentFlags().StringVar(&RedisUrl, "redis-url", "localhost:6379", "URL of the Redis instance to use")
	rootCmd.PersistentFlags().StringVar(&RedisPrefix, "redis-prefix", "godibs", "Redis channel prefix")
	rootCmd.PersistentFlags().StringVar(&ON, "on", OnNative, "System to run on (native|docker)")
	rootCmd.PersistentFlags().StringVar(&PLATFORM, "platform", PlatformDefault, "Platform specified in configuration to use (\""+PlatformAll+"\" to run for every platform)")
}

// Execute starts the main entry command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
	}

	if !(ON == OnNative || ON == OnDocker) {
		log.Fatal("Unsupported value for --on, must be either \""+OnNative+"\" or \""+OnDocker+"\"", rz.String("--on", ON))
	}
}
