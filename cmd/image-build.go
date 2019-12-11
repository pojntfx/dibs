package cmd

import "github.com/spf13/cobra"

var imageBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		if PLATFORM == PLATFORM_ALL {
			buildConfigs.BuildImageAll()
		} else {
			buildConfigs.BuildImage(PLATFORM)
		}
	},
}

func init() {
	imageCmd.AddCommand(imageBuildCmd)
}
