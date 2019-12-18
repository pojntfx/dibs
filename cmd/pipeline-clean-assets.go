package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineCleanAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Clean assets",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
		if err != nil {
			utils.LogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if err := platform.Assets.Clean(); err != nil {
				utils.LogErrorFatalPlatformSpecific("Could not clean assets", err, platform.Platform)
			}
		}
	},
}

func init() {
	PipelineCleanCmd.AddCommand(PipelineCleanAssetsCmd)
}
