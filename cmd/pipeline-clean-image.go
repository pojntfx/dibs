package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineCleanImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Clean image",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
		if err != nil {
			utils.LogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if output, err := platform.Assets.Build.CleanImage(platform.Platform); err != nil {
				utils.LogErrorFatalPlatformSpecific("Could not clean image", err, platform.Platform, output)
			}
		}
	},
}

func init() {
	PipelineCleanCmd.AddCommand(PipelineCleanImageCmd)
}
