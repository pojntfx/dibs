package v2

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineCleanAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Clean assets",
	Run: func(cmd *cobra.Command, args []string) {
		platforms, err := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if err := platform.Assets.Clean(); err != nil {
				utils.PipeLogErrorFatal("Could not clean assets", err, platform.Platform)
			}
		}
	},
}

func init() {
	PipelineCleanCmd.AddCommand(PipelineCleanAssetsCmd)
}
