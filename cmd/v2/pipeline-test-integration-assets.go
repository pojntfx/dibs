package v2

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineTestIntegrationAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Integration test the assets",
	Run: func(cmd *cobra.Command, args []string) {
		platforms, err := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if Executor == ExecutorDocker {
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
