package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelinePushImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Push image",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
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
