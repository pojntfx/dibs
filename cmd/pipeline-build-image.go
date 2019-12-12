package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineBuildImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Build image",
	Run: func(cmd *cobra.Command, args []string) {
		platforms, err := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if output, err := platform.Assets.Build.BuildImage(platform.Platform); err != nil {
				utils.PipeLogErrorFatal("Could not build image", err, platform.Platform, output)
			}
		}
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildImageCmd)
}
