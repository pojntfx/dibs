package cmd

import "github.com/spf13/cobra"

var binaryCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the binary output",
	Run: func(cmd *cobra.Command, args []string) {
		if PLATFORM == PLATFORM_ALL {
			switch ON {
			case ON_NATIVE:
				buildConfigs.BuildCleanAll()
			case ON_DOCKER:
				buildConfigs.BuildInDockerCleanAll()
			}
		} else {
			switch ON {
			case ON_NATIVE:
				buildConfigs.BuildClean(PLATFORM)
			case ON_DOCKER:
				buildConfigs.BuildInDockerClean(PLATFORM)
			}
		}
	},
}

func init() {
	binaryCmd.AddCommand(binaryCleanCmd)
}
