package cmd

import "github.com/spf13/cobra"

var manifestCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the Docker manifest output",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.Clean(PLATFORM)
	},
}

func init() {
	manifestCmd.AddCommand(manifestCleanCmd)
}
