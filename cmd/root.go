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
)

var rootCmd = &cobra.Command{
	Use:   "godibs",
	Short: "Distributed build system for Go",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&REDIS_URL, "redis-url", "localhost:6379", "Redis instance URL")
	rootCmd.PersistentFlags().StringVar(&REDIS_PREFIX, "redis-prefix", "godibs", "Redis channel prefix")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
	}
}
