package cmd

import "github.com/spf13/cobra"

var imagePushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push the Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		if PLATFORM == PLATFORM_ALL {
			buildConfigs.PushDockerImageAll()
		} else {
			buildConfigs.PushDockerImage(PLATFORM)
		}
	},
}

func init() {
	imageCmd.AddCommand(imagePushCmd)
}
