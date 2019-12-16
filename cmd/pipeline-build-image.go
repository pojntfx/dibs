package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineBuildImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Build image",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
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
