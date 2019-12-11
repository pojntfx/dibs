package cmd

import "github.com/spf13/cobra"

var imageIntegrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the Docker image integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.TestIntegrationDocker(PLATFORM)
	},
}

func init() {
	imageCmd.AddCommand(imageIntegrationtestCmd)
}
