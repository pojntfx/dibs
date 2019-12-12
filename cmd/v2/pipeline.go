package v2

import "github.com/spf13/cobra"

var PipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Pipeline building blocks",
}

func init() {
	RootCmd.AddCommand(PipelineCmd)
}
