package cmd

import "github.com/spf13/cobra"

// integrationtestCmd ist the command to start the server
var integrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.TestIntegrationGo(PLATFORM)
	},
}

// integrationtestCmd is a subcommand of langCmd
func init() {
	langCmd.AddCommand(integrationtestCmd)
}
