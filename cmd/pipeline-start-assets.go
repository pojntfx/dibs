package cmd

import (
	"github.com/pojntfx/dibs/pkg/pipes"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineStartAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Start assets",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		Dibs.RunForPlatformsConcurrently(platformFromConfig, platformFromConfig == PlatformAll, func(platform pipes.Platform) {
			if err := platform.Starters.Assets.StartStdoutStderr(platform.Platform); err != nil {
				utils.LogErrorFatalPlatformSpecific("Could not start", err, platform.Platform)
			}
		})
	},
}

func init() {
	PipelineStartCmd.AddCommand(PipelineStartAssetsCmd)
}
