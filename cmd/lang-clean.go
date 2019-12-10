package cmd

import "github.com/spf13/cobra"

// cleanCmd ist the command to start the server
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the language output",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.Clean(PLATFORM)
	},
}

// cleanCmd is a subcommand of langCmd
func init() {
	langCmd.AddCommand(cleanCmd)
}
