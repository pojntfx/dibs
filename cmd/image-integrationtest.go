package cmd

import "github.com/spf13/cobra"

var imageIntegrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the Docker image integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		switch ON {
		case ON_NATIVE:
			buildConfigs.TestIntegrationImage(PLATFORM)
		case ON_DOCKER:
			buildConfigs.TestIntegrationImageInDocker(PLATFORM)
		}
	},
}

func init() {
	imageCmd.AddCommand(imageIntegrationtestCmd)
}
