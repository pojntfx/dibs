package v2

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineTestIntegrationImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Integration test the image",
	Run: func(cmd *cobra.Command, args []string) {
		platforms, err := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if Executor == ExecutorDocker {
				if output, err := platform.Tests.Integration.Image.BuildImage(platform.Platform); err != nil {
					utils.PipeLogErrorFatal("Could not build image integration test image", err, platform.Platform, output)
				}
				output, err := platform.Tests.Integration.Image.StartImage(platform.Platform)
				utils.PipeLogErrorInfo("Image integration test ran in Docker", err, platform.Platform, output)
			} else {
				output, err := platform.Tests.Integration.Image.Start(platform.Platform)
				utils.PipeLogErrorInfo("Image integration test ran", err, platform.Platform, output)
			}
		}
	},
}

func init() {
	PipelineTestIntegrationCmd.AddCommand(PipelineTestIntegrationImageCmd)
}
