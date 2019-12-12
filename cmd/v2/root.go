package v2

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var rootCmd = &cobra.Command{
	Use:   "dibs",
	Short: "System for distributed polyglot, multi-module and multi-architecture development",
}

var (
	Platform string

	Executor string

	RedisUrl    string
	RedisPrefix string

	Dibs utils.Dibs
)

const (
	PlatformAll     = "all"
	PlatformDefault = PlatformAll

	ExecutorNative  = "native"
	ExecutorDocker  = "docker"
	ExecutorDefault = ExecutorNative

	RedisUrlDefault    = "localhost:6379"
	RedisPrefixDefault = "dibs"

	DibsPath = "."
	DibsFile = ".dibs"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&Platform, "platform", "p", PlatformDefault, `Platform to run on ("`+PlatformAll+`" runs on all platforms specified in configuration file)`)
	rootCmd.PersistentFlags().StringVarP(&Executor, "executor", "e", ExecutorDefault, `Executor to run on `+`("`+ExecutorDocker+`"|"`+ExecutorNative+`")`)
	rootCmd.PersistentFlags().StringVarP(&RedisUrl, "redis-url", "u", RedisUrlDefault, "URL of the Redis instance to use")
	rootCmd.PersistentFlags().StringVarP(&RedisPrefix, "redis-prefix", "c", RedisPrefixDefault, "Redis channel prefix to use")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Could not start root command", rz.Err(err))
	}

	if !(Executor == ExecutorDocker || Executor == ExecutorNative) {
		log.Fatal("Unsupported value for --executor, must be either \""+ExecutorDocker+"\" or \""+ExecutorNative+"\"", rz.String("--executor", Executor))
	}

	viper.AddConfigPath(DibsPath)
	viper.SetConfigName(DibsFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Could not read config", rz.Err(err))
	}

	if err := viper.Unmarshal(&Dibs); err != nil {
		log.Fatal("Could not unmarshal config", rz.Err(err))
	}
}
