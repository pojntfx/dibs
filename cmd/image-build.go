package cmd

import "github.com/spf13/cobra"

var imageBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		switch ON {
		case ON_NATIVE:
			buildConfigs.BuildDocker(PLATFORM)
		case ON_DOCKER:
			buildConfigs.BuildDockerInDocker(PLATFORM)
		}
	},
}

func init() {
	imageCmd.AddCommand(imageBuildCmd)
}
