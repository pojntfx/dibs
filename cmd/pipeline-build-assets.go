package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineBuildAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Build assets",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
		if err != nil {
			utils.LogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if viper.GetString(ExecutorKey) == ExecutorDocker {
				if output, err := platform.Assets.GetAssetsFromDockerImage(platform.Platform); err != nil {
					utils.LogErrorFatalPlatformSpecific("Could not get assets from Docker image", err, platform.Platform, output)
				}
			} else {
				if output, err := platform.Assets.Build.Start(platform.Platform); err != nil {
					utils.LogErrorFatalPlatformSpecific("Could not build assets", err, platform.Platform, output)
				}
			}
		}
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildAssetsCmd)
}
