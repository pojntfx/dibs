package cmd

import (
	"github.com/pojntfx/dibs/pkg/pipes"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineCleanAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Clean assets",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		Dibs.RunForPlatformsConcurrently(platformFromConfig, platformFromConfig == PlatformAll, func(platform pipes.Platform) {
			if err := platform.Assets.Clean(); err != nil {
				utils.LogErrorFatalPlatformSpecific("Could not clean assets", err, platform.Platform)
			}
		})
	},
}

func init() {
	PipelineCleanCmd.AddCommand(PipelineCleanAssetsCmd)
}
