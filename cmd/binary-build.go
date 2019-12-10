package cmd

import "github.com/spf13/cobra"

var binaryBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the binary",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.Build(PLATFORM)
	},
}

func init() {
	binaryCmd.AddCommand(binaryBuildCmd)
}
