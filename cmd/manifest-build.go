package cmd

import "github.com/spf13/cobra"

var manifestBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Docker manifest",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.BuildDockerManifest()
	},
}

func init() {
	manifestCmd.AddCommand(manifestBuildCmd)
}
