package cmd

import "github.com/spf13/cobra"

var PipelineCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the output of a pipeline building block",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ReadConfig()
	},
}

func init() {
	PipelineCmd.AddCommand(PipelineCleanCmd)
}
