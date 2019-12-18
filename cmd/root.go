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
	Short: "System for distributed polyglot, multi-module, multi-architecture development and CI/CD",
	Long: `System for distributed polyglot, multi-module, multi-architecture development and CI/CD

For full functionality, it requires the following binaries to be in PATH:

- "docker"	(see https://www.docker.com/)
- "kubectl"	(see https://kubernetes.io/docs/reference/kubectl/)
- "helm"	(see https://helm.sh/)
- "skaffold"	(see https://skaffold.dev/)
- "ghr"		(see https://github.com/tcnksm/ghr)
- "cr"		(see https://github.com/helm/chart-releaser)`,
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

	if err := viper.BindPFlag(PlatformKey, RootCmd.PersistentFlags().Lookup(platformFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}
	if err := viper.BindPFlag(ExecutorKey, RootCmd.PersistentFlags().Lookup(executorFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindPFlag(DibsFileKey, RootCmd.PersistentFlags().Lookup(dibsFileFlag)); err != nil {
		log.Fatal("Could not bind flag", rz.Err(err))
	}

	if err := viper.BindEnv(PlatformKey, PlatformEnvDocker); err != nil {
		log.Fatal("Could not bind key", rz.Err(err))
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
