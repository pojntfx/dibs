package cmd

import "github.com/spf13/cobra"

var langUnittestCmd = &cobra.Command{
	Use:   "unittest",
	Short: "Run the unit tests",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.TestUnit(PLATFORM)
	},
}

func init() {
	langCmd.AddCommand(langUnittestCmd)
}
