package cmd

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelinePushImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Push image",
	Run: func(cmd *cobra.Command, args []string) {
		platforms, err := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if output, err := platform.Assets.Build.PushImage(platform.Platform); err != nil {
				utils.PipeLogErrorFatal("Could not push image", err, platform.Platform, output)
			}
		}
	},
}

func init() {
	PipelinePushCmd.AddCommand(PipelinePushImageCmd)
}
