package cmd

import (
	"errors"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var PipelineSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Develop with multiple modules",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		lang := viper.GetString(LangKey)

		if !(lang == LangGo) {
			return errors.New(`unsupported value "` + lang + `" for --lang, must be "` + LangGo + `"`)
		}

		return nil
	},
}

func init() {
	var (
		lang string

		redisUrl      string
		redisPrefix   string
		redisPassword string

		langFlag = strings.Replace(LangKey, "_", "-", -1)

		redisUrlFlag      = strings.Replace(strings.Replace(SyncRedisUrlKey, SyncKeyPrefix, "", -1), "_", "-", -1)
		redisPrefixFlag   = strings.Replace(strings.Replace(SyncRedisPrefixKey, SyncKeyPrefix, "", -1), "_", "-", -1)
		redisPasswordFlag = strings.Replace(strings.Replace(SyncRedisPasswordKey, SyncKeyPrefix, "", -1), "_", "-", -1)
	)

	PipelineSyncCmd.PersistentFlags().StringVarP(&lang, langFlag, "l", LangDefault, `Language to develop the modules for (currently only "`+LangGo+`" is supported)`)

	PipelineSyncCmd.PersistentFlags().StringVarP(&redisUrl, redisUrlFlag, "u", SyncClientRedisUrlDefault, "URL of the Redis instance to use")
	PipelineSyncCmd.PersistentFlags().StringVarP(&redisPrefix, redisPrefixFlag, "c", SyncClientRedisPrefixDefault, "Redis channel prefix to use")
	PipelineSyncCmd.PersistentFlags().StringVarP(&redisPassword, redisPasswordFlag, "s", SyncClientRedisPasswordDefault, "Redis password to use")

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(LangKey, PipelineSyncCmd.PersistentFlags().Lookup(langFlag)); err != nil {
		utils.LogErrorCouldNotBindFlag(err)
	}

	if err := viper.BindPFlag(SyncRedisUrlKey, PipelineSyncCmd.PersistentFlags().Lookup(redisUrlFlag)); err != nil {
		utils.LogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(SyncRedisPrefixKey, PipelineSyncCmd.PersistentFlags().Lookup(redisPrefixFlag)); err != nil {
		utils.LogErrorCouldNotBindFlag(err)
	}
	if err := viper.BindPFlag(SyncRedisPasswordKey, PipelineSyncCmd.PersistentFlags().Lookup(redisPasswordFlag)); err != nil {
		utils.LogErrorCouldNotBindFlag(err)
	}

	viper.AutomaticEnv()

	PipelineCmd.AddCommand(PipelineSyncCmd)
}
