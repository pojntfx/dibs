package cmd

import "github.com/spf13/cobra"

var PipelineTestUnitCmd = &cobra.Command{
	Use:   "unit",
	Short: "Unit test artifacts",
}

func init() {
	PipelineTestCmd.AddCommand(PipelineTestUnitCmd)
}
