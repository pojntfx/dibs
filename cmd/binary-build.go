package cmd

import "github.com/spf13/cobra"

var binaryBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the binary",
	Run: func(cmd *cobra.Command, args []string) {
		switch ON {
		case ON_NATIVE:
			buildConfigs.Build(PLATFORM)
		case ON_DOCKER:
			buildConfigs.BuildInDocker(PLATFORM)
		}
	},
}

func init() {
	binaryCmd.AddCommand(binaryBuildCmd)
}