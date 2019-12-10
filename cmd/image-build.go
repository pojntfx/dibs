package cmd

import "github.com/spf13/cobra"

var imageBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.BuildDocker(PLATFORM)
	},
}

func init() {
	imageCmd.AddCommand(imageBuildCmd)
}
