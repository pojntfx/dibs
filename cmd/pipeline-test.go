package cmd

import "github.com/spf13/cobra"

var PipelineTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test with a pipeline building block",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ReadConfig()
	},
}

func init() {

	PipelineCmd.AddCommand(PipelineTestCmd)
}
