package cmd

import (
	"github.com/pojntfx/dibs/pkg/pipes"
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineTestIntegrationChartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Integration test the chart",
	Run: func(cmd *cobra.Command, args []string) {
		platformFromConfig := viper.GetString(PlatformKey)

		Dibs.RunForPlatformsConcurrently(platformFromConfig, platformFromConfig == PlatformAll, func(platform pipes.Platform) {
			if viper.GetString(ExecutorKey) == ExecutorDocker {
				if output, err := platform.Tests.Integration.Chart.BuildImage(platform.Platform); err != nil {
					utils.LogErrorFatalPlatformSpecific("Could not build chart integration test chart", err, platform.Platform, output)
				}
				output, err := platform.Tests.Integration.Chart.StartImage(platform.Platform)
				utils.LogErrorInfo("Chart integration test ran in Docker", err, platform.Platform, output)
			} else {
				output, err := platform.Tests.Integration.Chart.Start(platform.Platform)
				utils.LogErrorInfo("Chart integration test ran", err, platform.Platform, output)
			}
		})
	},
}

func init() {
	PipelineTestIntegrationCmd.AddCommand(PipelineTestIntegrationChartCmd)
}
