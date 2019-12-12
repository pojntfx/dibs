package v2

import "github.com/spf13/cobra"

var PipelineTestIntegrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Integration test with a pipeline building block",
}

func init() {
	PipelineTestCmd.AddCommand(PipelineTestIntegrationCmd)
}
