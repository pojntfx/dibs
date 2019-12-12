package cmd

import "github.com/spf13/cobra"

var PipelineTestUnitCmd = &cobra.Command{
	Use:   "unit",
	Short: "Unit test with a pipeline building block",
}

func init() {
	PipelineTestCmd.AddCommand(PipelineTestUnitCmd)
}
