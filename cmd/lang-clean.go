package cmd

import "github.com/spf13/cobra"

var langCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the language output",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.Clean(PLATFORM)
	},
}

func init() {
	langCmd.AddCommand(langCleanCmd)
}
