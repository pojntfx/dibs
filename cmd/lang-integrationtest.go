package cmd

import "github.com/spf13/cobra"

var langIntegrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		switch ON {
		case ON_NATIVE:
			buildConfigs.TestIntegrationGo(PLATFORM)
		case ON_DOCKER:
			buildConfigs.TestIntegrationGoInDocker(PLATFORM)
		}
	},
}

func init() {
	langCmd.AddCommand(langIntegrationtestCmd)
}
