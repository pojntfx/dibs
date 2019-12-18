package cmd

import "github.com/spf13/cobra"

var PipelineTestIntegrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Integration test artifacts",
}

func init() {
	PipelineTestCmd.AddCommand(PipelineTestIntegrationCmd)
}
