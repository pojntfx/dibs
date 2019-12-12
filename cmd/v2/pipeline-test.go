package v2

import "github.com/spf13/cobra"

var PipelineTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test with a pipeline building block",
}

func init() {
	PipelineCmd.AddCommand(PipelineTestCmd)
}
