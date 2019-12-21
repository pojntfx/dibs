package cmd

import (
	"github.com/pojntfx/dibs/pkg/pipes"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineTestIntegrationImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Integration test the image",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		Dibs.RunForPlatformsConcurrently(platformFromConfig, platformFromConfig == PlatformAll, func(platform pipes.Platform) {
			if viper.GetString(ExecutorKey) == ExecutorDocker {
				if output, err := platform.Tests.Integration.Image.BuildImage(platform.Platform); err != nil {
					utils.LogErrorFatalPlatformSpecific("Could not build image integration test image", err, platform.Platform, output)
				}
				output, err := platform.Tests.Integration.Image.StartImage(platform.Platform)
				utils.LogErrorInfo("Image integration test ran in Docker", err, platform.Platform, output)
			} else {
				output, err := platform.Tests.Integration.Image.Start(platform.Platform)
				utils.LogErrorInfo("Image integration test ran", err, platform.Platform, output)
			}
		})
	},
}

func init() {
	PipelineTestIntegrationCmd.AddCommand(PipelineTestIntegrationImageCmd)
}
