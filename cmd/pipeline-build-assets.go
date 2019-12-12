package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineBuildAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Build assets",
	Run: func(cmd *cobra.Command, args []string) {
		platforms, err := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if Executor == ExecutorDocker {
				if output, err := platform.Assets.GetAssetsFromDockerImage(platform.Platform); err != nil {
					utils.PipeLogErrorFatal("Could not get assets from Docker image", err, platform.Platform, output)
				}
			} else {
				if output, err := platform.Assets.Build.Start(platform.Platform); err != nil {
					utils.PipeLogErrorFatal("Could not build assets", err, platform.Platform, output)
				}
			}
		}
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildAssetsCmd)
}
