package cmd

import (
	"errors"
	"github.com/pojntfx/dibs/pkg/pipes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
	"strings"
)

var RootCmd = &cobra.Command{
	Use:   "dibs",
	Short: "System for distributed polyglot, multi-module and multi-architecture development",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		executor := viper.GetString(ExecutorKey)

		if !(executor == ExecutorDocker || executor == ExecutorNative) {
			return errors.New(`unsupported value "` + executor + `" for --executor, must be either "` + ExecutorDocker + `" or "` + ExecutorNative + `"`)
		}

		return nil
	},
}

var (
	Dibs pipes.Dibs
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

	EnvPrefix = "dibs"

	PlatformKey = "platform"
	ExecutorKey = "executor"

	DibsFileKey = "config_file"
)

func init() {
	var (
		platform string

		executor string

		dibsFile = DibsName + ".yml"

		platformFlag = strings.Replace(PlatformKey, "_", "-", -1)
		executorFlag = strings.Replace(ExecutorKey, "_", "-", -1)

		dibsFileFlag = strings.Replace(DibsFileKey, "_", "-", -1)
	)

	RootCmd.PersistentFlags().StringVarP(&platform, platformFlag, "p", PlatformDefault, `Platform to run on ("`+PlatformAll+`" runs on all platforms specified in configuration file)`)
	RootCmd.PersistentFlags().StringVarP(&executor, executorFlag, "e", ExecutorDefault, `Executor to run on `+`("`+ExecutorDocker+`"|"`+ExecutorNative+`")`)

	RootCmd.PersistentFlags().StringVarP(&dibsFile, dibsFileFlag, "f", DibsFileDefault, "Configuration file to use")

	viper.SetEnvPrefix(EnvPrefix)

	if err := viper.BindPFlag(PlatformKey, PipelineSyncClientCmd.PersistentFlags().Lookup(platformFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(ExecutorKey, PipelineSyncClientCmd.PersistentFlags().Lookup(executorFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(DibsFileKey, PipelineSyncClientCmd.PersistentFlags().Lookup(dibsFileFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	viper.AutomaticEnv()
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal("Could not start root command", rz.Err(err))
	}
}

func ReadConfig(dibsFile string) error {
	viper.AddConfigPath(DibsPath)

	if dibsFile != DibsFileDefault {
		viper.SetConfigFile(dibsFile)
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
