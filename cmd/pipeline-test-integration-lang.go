package cmd

import (
	"github.com/pojntfx/dibs/pkg/pipes"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineTestIntegrationLangCmd = &cobra.Command{
	Use:   "lang",
	Short: "Integration test with the language's toolchain",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		Dibs.RunForPlatformsConcurrently(platformFromConfig, platformFromConfig == PlatformAll, func(platform pipes.Platform) {
			if viper.GetString(ExecutorKey) == ExecutorDocker {
				if output, err := platform.Tests.Integration.Lang.BuildImage(platform.Platform); err != nil {
					utils.LogErrorFatalPlatformSpecific("Could not build lang integration test image", err, platform.Platform, output)
				}
				output, err := platform.Tests.Integration.Lang.StartImage(platform.Platform)
				utils.LogErrorInfo("Lang integration test ran in Docker", err, platform.Platform, output)
			} else {
				output, err := platform.Tests.Integration.Lang.Start(platform.Platform)
				utils.LogErrorInfo("Lang integration test ran", err, platform.Platform, output)
			}
		})
	},
}

func init() {
	PipelineTestIntegrationCmd.AddCommand(PipelineTestIntegrationLangCmd)
}
