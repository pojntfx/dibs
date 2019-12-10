package cmd

import "github.com/spf13/cobra"

var binaryIntegrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.TestIntegrationGo(PLATFORM)
	},
}

func init() {
	binaryCmd.AddCommand(binaryIntegrationtestCmd)
}
