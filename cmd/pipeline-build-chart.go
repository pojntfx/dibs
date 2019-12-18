package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineBuildChartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Build chart",
	Run: func(cmd *cobra.Command, args []string) {
		if output, err := Dibs.BuildHelmChart(viper.GetString(PlatformKey)); err != nil {
			utils.LogErrorFatal("Could not build chart", err, output)
		}
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildChartCmd)
}
