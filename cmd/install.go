package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and start the project",
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
			if err := platform.Assets.Build.StartChart(platform.Platform, platform.ChartProfiles.Production); err != nil {
				utils.PipeLogErrorFatalWithProfile("Could not install profile", err, platform.Platform, platform.ChartProfiles.Production)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(InstallCmd)
}
