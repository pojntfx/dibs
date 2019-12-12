package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineCleanImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Clean image",
	Run: func(cmd *cobra.Command, args []string) {
		platforms, err := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if output, err := platform.Assets.Build.CleanImage(platform.Platform); err != nil {
				utils.PipeLogErrorFatal("Could not clean image", err, platform.Platform, output)
			}
		}
	},
}

func init() {
	PipelineCleanCmd.AddCommand(PipelineCleanImageCmd)
}
