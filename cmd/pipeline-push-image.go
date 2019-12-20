package cmd

import (
	"github.com/pojntfx/dibs/pkg/pipes"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelinePushImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Push image",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		Dibs.RunForPlatforms(platformFromConfig, platformFromConfig == PlatformAll, func(platform pipes.Platform) {
			if output, err := platform.Assets.Build.PushImage(platform.Platform); err != nil {
				utils.LogErrorFatalPlatformSpecific("Could not push image", err, platform.Platform, output)
			}
		})
	},
}

func init() {
	PipelinePushCmd.AddCommand(PipelinePushImageCmd)
}
