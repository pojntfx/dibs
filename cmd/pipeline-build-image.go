package cmd

import (
	"github.com/pojntfx/dibs/pkg/pipes"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineBuildImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Build image",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		Dibs.RunForPlatformsConcurrently(platformFromConfig, platformFromConfig == PlatformAll, func(platform pipes.Platform) {
			if output, err := platform.Assets.Build.BuildImage(platform.Platform); err != nil {
				utils.LogErrorFatalPlatformSpecific("Could not build image", err, platform.Platform, output)
			}
		})
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildImageCmd)
}
