package cmd

import (
	"errors"
	"github.com/pojntfx/dibs/pkg/pipes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var RootCmd = &cobra.Command{
	Use:   "dibs",
	Short: "System for distributed polyglot, multi-module and multi-architecture development",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !(Executor == ExecutorDocker || Executor == ExecutorNative) {
			return errors.New(`unsupported value "` + Executor + `" for --executor, must be either "` + ExecutorDocker + `" or "` + ExecutorNative + `"`)
		}

		return nil
	},
}

var (
	Platform string

	Executor string

	RedisUrl    string
	RedisPrefix string

	DibsFile = DibsName + ".yml"
	Dibs     pipes.Dibs
)

const (
	PlatformAll     = "all"
	PlatformDefault = PlatformAll

	ExecutorNative  = "native"
	ExecutorDocker  = "docker"
	ExecutorDefault = ExecutorNative

	DibsPath        = "."
	DibsName        = ".dibs"
	DibsFileDefault = DibsName + ".yml"
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&Platform, "platform", "p", PlatformDefault, `Platform to run on ("`+PlatformAll+`" runs on all platforms specified in configuration file)`)
	RootCmd.PersistentFlags().StringVarP(&Executor, "executor", "e", ExecutorDefault, `Executor to run on `+`("`+ExecutorDocker+`"|"`+ExecutorNative+`")`)

	RootCmd.PersistentFlags().StringVarP(&DibsFile, "config-file", "f", DibsFileDefault, "Configuration file to use")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal("Could not start root command", rz.Err(err))
	}
}

func ReadConfig() error {
	viper.AddConfigPath(DibsPath)

	if DibsFile != DibsFileDefault {
		viper.SetConfigFile(DibsFile)
	} else {
		viper.SetConfigName(DibsName)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&Dibs); err != nil {
		return err
	}

	return nil
}
