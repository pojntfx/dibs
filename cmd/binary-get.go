package cmd

import "github.com/spf13/cobra"

var binaryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the binary from the Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.GetBinaryFromDockerImage(PLATFORM)
	},
}

func init() {
	binaryCmd.AddCommand(binaryGetCmd)
}