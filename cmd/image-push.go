package cmd

import "github.com/spf13/cobra"

var imagePushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push the Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.PushDockerImage(PLATFORM)
	},
}

func init() {
	imageCmd.AddCommand(imagePushCmd)
}
