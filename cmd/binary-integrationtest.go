package cmd

import "github.com/spf13/cobra"

var binaryIntegrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the binary integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		switch ON {
		case ON_NATIVE:
			buildConfigs.TestIntegrationBinary(PLATFORM)
		case ON_DOCKER:
			buildConfigs.TestIntegrationBinaryInDocker(PLATFORM)
		}
	},
}

func init() {
	binaryCmd.AddCommand(binaryIntegrationtestCmd)
}
