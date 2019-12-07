package cmd

import (
	"github.com/spf13/cobra"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var (
	REDIS_URL    string
	REDIS_PREFIX string
)

const (
	REDIS_SUFFIX_UP_BUILT        = "up_built"
	REDIS_SUFFIX_UP_TESTED       = "up_tested"
	REDIS_SUFFIX_UP_STARTED      = "up_started"
	REDIS_SUFFIX_UP_REGISTERED   = "up_registered"
	REDIS_SUFFIX_UP_UNREGISTERED = "up_unregistered"
	REDIS_SUFFIX_UP_PUSHED       = "up_pushed"
	REDIS_SUFFIX_DOWN_DOWNLOADED = "down_downloaded"
)

// serverCmd ist the main entry command
var rootCmd = &cobra.Command{
	Use:   "godibs",
	Short: "System for distributed multi-module development with Go",
}

// init maps the flags to the config
func init() {
	rootCmd.PersistentFlags().StringVar(&REDIS_URL, "redis-url", "localhost:6379", "URL of the Redis instance to use")
	rootCmd.PersistentFlags().StringVar(&REDIS_PREFIX, "redis-prefix", "godibs", "Redis channel prefix")
}

// Execute starts the main entry command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
	}
}
