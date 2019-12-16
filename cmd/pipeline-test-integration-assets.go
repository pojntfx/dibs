package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineTestIntegrationAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Integration test the assets",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if viper.GetString(ExecutorKey) == ExecutorDocker {
				if output, err := platform.Tests.Integration.Assets.BuildImage(platform.Platform); err != nil {
					utils.PipeLogErrorFatal("Could not build assets integration test image", err, platform.Platform, output)
				}
				output, err := platform.Tests.Integration.Assets.StartImage(platform.Platform)
				utils.PipeLogErrorInfo("Assets integration test ran in Docker", err, platform.Platform, output)
			} else {
				output, err := platform.Tests.Integration.Assets.Start(platform.Platform)
				utils.PipeLogErrorInfo("Assets integration test ran", err, platform.Platform, output)
			}
		}
	},
}

func init() {
	PipelineTestIntegrationCmd.AddCommand(PipelineTestIntegrationAssetsCmd)
}
