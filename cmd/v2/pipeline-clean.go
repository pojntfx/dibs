package v2

import "github.com/spf13/cobra"

var PipelineCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the output of a pipeline building block",
}

func init() {
	PipelineCmd.AddCommand(PipelineCleanCmd)
}
