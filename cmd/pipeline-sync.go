package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"strings"
)

var PipelineSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with a pipeline building block",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		lang := viper.GetString(LangKey)

		if !(lang == LangGo) {
			return errors.New(`unsupported value "` + lang + `" for --lang, must be "` + LangGo + `"`)
		}

		return nil
	},
}

const (
	LangDefault = LangGo
	LangGo      = "go"

	RedisUrlDefault    = "localhost:6379"
	RedisPrefixDefault = "dibs"

	LangKey = "lang"

	RedisUrlKey    = "redis_url"
	RedisPrefixKey = "redis_prefix"
)

func init() {
	var (
		lang string

		redisUrl    string
		redisPrefix string

		langFlag = strings.Replace(LangKey, "_", "-", -1)

		redisUrlFlag    = strings.Replace(RedisUrlKey, "_", "-", -1)
		redisPrefixFlag = strings.Replace(RedisPrefixKey, "_", "-", -1)
	)

	PipelineSyncCmd.PersistentFlags().StringVarP(&lang, langFlag, "l", LangDefault, `Language to develop the modules for (currently only "`+LangGo+`" is supported)`)

	PipelineSyncCmd.PersistentFlags().StringVarP(&redisUrl, redisUrlFlag, "u", RedisUrlDefault, "URL of the Redis instance to use")
	PipelineSyncCmd.PersistentFlags().StringVarP(&redisPrefix, redisPrefixFlag, "c", RedisPrefixDefault, "Redis channel prefix to use")

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(LangKey, PipelineSyncCmd.PersistentFlags().Lookup(langFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(RedisUrlKey, PipelineSyncCmd.PersistentFlags().Lookup(redisUrlFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(RedisPrefixKey, PipelineSyncCmd.PersistentFlags().Lookup(redisPrefixFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	viper.AutomaticEnv()

	PipelineCmd.AddCommand(PipelineSyncCmd)
}
