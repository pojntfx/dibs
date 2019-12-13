package cmd

import "github.com/spf13/cobra"

var PipelineBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build with a pipeline building block",
}

func init() {
	PipelineCmd.AddCommand(PipelineBuildCmd)
}