package cmd

import "github.com/spf13/cobra"

var binaryCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the binary output",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.Clean(PLATFORM)
	},
}

func init() {
	binaryCmd.AddCommand(binaryCleanCmd)
}
