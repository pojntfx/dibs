package cmd

import "github.com/spf13/cobra"

var manifestPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push the Docker manifest",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.PushDockerManifest()
	},
}

func init() {
	manifestCmd.AddCommand(manifestPushCmd)
}
