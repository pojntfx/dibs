package cmd

import "github.com/spf13/cobra"

var PipelinePushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push with a pipeline building block",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ReadConfig()
	},
}

func init() {
	PipelineCmd.AddCommand(PipelinePushCmd)
}
