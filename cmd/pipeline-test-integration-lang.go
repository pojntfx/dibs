package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineTestIntegrationLangCmd = &cobra.Command{
	Use:   "lang",
	Short: "Integration test with the language's toolchain",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if viper.GetString(ExecutorKey) == ExecutorDocker {
				if output, err := platform.Tests.Integration.Lang.BuildImage(platform.Platform); err != nil {
					utils.PipeLogErrorFatal("Could not build lang integration test image", err, platform.Platform, output)
				}
				output, err := platform.Tests.Integration.Lang.StartImage(platform.Platform)
				utils.PipeLogErrorInfo("Lang integration test ran in Docker", err, platform.Platform, output)
			} else {
				output, err := platform.Tests.Integration.Lang.Start(platform.Platform)
				utils.PipeLogErrorInfo("Lang integration test ran", err, platform.Platform, output)
			}
		}
	},
}

func init() {
	PipelineTestIntegrationCmd.AddCommand(PipelineTestIntegrationLangCmd)
}
