package cmd

import "github.com/spf13/cobra"

var PipelineSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with a pipeline building block",
}

func init() {
	PipelineCmd.AddCommand(PipelineSyncCmd)
}
