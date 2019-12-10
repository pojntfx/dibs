package cmd

import "github.com/spf13/cobra"

// unittestCmd ist the command to start the server
var unittestCmd = &cobra.Command{
	Use:   "unittest",
	Short: "Run the unit tests",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.TestUnit(PLATFORM)
	},
}

// unittestCmd is a subcommand of langCmd
func init() {
	langCmd.AddCommand(unittestCmd)
}
