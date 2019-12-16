package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Develop the project",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ReadConfig(viper.GetString(DibsFileKey))
	},
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if err := platform.Assets.Build.DevChart(platform.Platform, platform.ChartProfiles.Development); err != nil {
				utils.PipeLogErrorFatalWithProfile("Could not dev on profile", err, platform.Platform, platform.ChartProfiles.Development)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(DevCmd)
}
