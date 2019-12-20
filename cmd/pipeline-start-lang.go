package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineStartLangCmd = &cobra.Command{
	Use:   "lang",
	Short: "Start with the language's toolchain",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		platforms, err := Dibs.GetPlatforms(platformFromConfig, platformFromConfig == PlatformAll)
		if err != nil {
			utils.LogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if err := platform.Starters.Lang.StartStdoutStderr(platform.Platform); err != nil {
				utils.LogErrorFatalPlatformSpecific("Could not start", err, platform.Platform)
			}
		}
	},
}

func init() {
	PipelineStartCmd.AddCommand(PipelineStartLangCmd)
}
