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
