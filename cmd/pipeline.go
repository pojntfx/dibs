package cmd

import "github.com/spf13/cobra"

var PipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Individual pipeline building blocks",
}

func init() {
	RootCmd.AddCommand(PipelineCmd)
}
