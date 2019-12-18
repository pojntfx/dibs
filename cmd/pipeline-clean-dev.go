package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineCleanDevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Clean development resources manually",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ReadConfig(viper.GetString(DibsFileKey))
	},
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
		if err != nil {
			utils.LogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if err := platform.Assets.Build.CleanStartedChart(platform.Platform, platform.ChartProfiles.Development); err != nil {
				utils.LogErrorFatalWithProfile("Could not uninstall profile", err, platform.Platform, platform.ChartProfiles.Development)
			}
		}
	},
}

func init() {
	PipelineCleanCmd.AddCommand(PipelineCleanDevCmd)
}
