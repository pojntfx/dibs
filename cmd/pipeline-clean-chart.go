package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineCleanChartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Clean chart",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Dibs.CleanHelmChart(); err != nil {
			utils.PipeLogErrorFatalNonPlatformSpecific("Could not clean chart", err)
		}
	},
}

func init() {
	PipelineCleanCmd.AddCommand(PipelineCleanChartCmd)
}
