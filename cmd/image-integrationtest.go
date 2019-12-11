package cmd

import "github.com/spf13/cobra"

var imageIntegrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the Docker image integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		switch ON {
		case ON_NATIVE:
			buildConfigs.TestIntegrationDocker(PLATFORM)
		case ON_DOCKER:
			buildConfigs.TestIntegrationDockerInDocker(PLATFORM)
		}
	},
}

func init() {
	imageCmd.AddCommand(imageIntegrationtestCmd)
}
