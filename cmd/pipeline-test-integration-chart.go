package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineTestIntegrationChartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Integration test the chart",
	Run: func(cmd *cobra.Command, args []string) {
		platforms, err := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if Executor == ExecutorDocker {
				if output, err := platform.Tests.Integration.Chart.BuildImage(platform.Platform); err != nil {
					utils.PipeLogErrorFatal("Could not build chart integration test chart", err, platform.Platform, output)
				}
				output, err := platform.Tests.Integration.Chart.StartImage(platform.Platform, struct {
					Key   string
					Value string
				}{
					Key:   "IP",
					Value: "192.168.178.53",
				})
				utils.PipeLogErrorInfo("Chart integration test ran in Docker", err, platform.Platform, output)
			} else {
				output, err := platform.Tests.Integration.Chart.Start(platform.Platform)
				utils.PipeLogErrorInfo("Chart integration test ran", err, platform.Platform, output)
			}
		}
	},
}

func init() {
	PipelineTestIntegrationCmd.AddCommand(PipelineTestIntegrationChartCmd)
}
