package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineBuildChartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Build chart",
	Run: func(cmd *cobra.Command, args []string) {
		if output, err := Dibs.BuildHelmChart(Platform); err != nil {
			utils.PipeLogErrorFatalNonPlatformSpecific("Could not build chart", err, output)
		}
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildChartCmd)
}
