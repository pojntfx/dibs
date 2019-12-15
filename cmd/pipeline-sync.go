package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var PipelineSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with a pipeline building block",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !(Lang == LangGo) {
			return errors.New(`unsupported value "` + Lang + `" for --lang, must be "` + LangGo + `"`)
		}

		return nil
	},
}

var (
	Lang string
)

const (
	LangDefault = LangGo
	LangGo      = "go"

	RedisUrlDefault    = "localhost:6379"
	RedisPrefixDefault = "dibs"
)

func init() {
	PipelineSyncCmd.PersistentFlags().StringVarP(&Lang, "lang", "l", LangDefault, `Language to develop the modules for (currently only "`+LangGo+`" is supported)`)

	PipelineSyncCmd.PersistentFlags().StringVarP(&RedisUrl, "redis-url", "u", RedisUrlDefault, "URL of the Redis instance to use")
	PipelineSyncCmd.PersistentFlags().StringVarP(&RedisPrefix, "redis-prefix", "c", RedisPrefixDefault, "Redis channel prefix to use")

	PipelineCmd.AddCommand(PipelineSyncCmd)
}
