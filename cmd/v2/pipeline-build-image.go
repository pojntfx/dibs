package v2

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineBuildImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Build image",
	Run: func(cmd *cobra.Command, args []string) {
		platforms, err := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
		if err != nil {
			utils.PipeLogErrorPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if output, err := platform.Binary.Build.BuildImage(platform.Platform); err != nil {
				utils.PipeLogError("Could not build image", err, platform.Platform, output)
			}
			return
		}
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildImageCmd)
}
