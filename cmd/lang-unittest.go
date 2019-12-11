package cmd

import "github.com/spf13/cobra"

var langUnittestCmd = &cobra.Command{
	Use:   "unittest",
	Short: "Run the unit tests",
	Run: func(cmd *cobra.Command, args []string) {
		if PLATFORM == PLATFORM_ALL {
			switch ON {
			case ON_NATIVE:
				buildConfigs.TestUnitAll()
			case ON_DOCKER:
				buildConfigs.TestUnitInDockerAll()
			}
		} else {
			switch ON {
			case ON_NATIVE:
				buildConfigs.TestUnit(PLATFORM)
			case ON_DOCKER:
				buildConfigs.TestUnitInDocker(PLATFORM)
			}
		}
	},
}

func init() {
	langCmd.AddCommand(langUnittestCmd)
}
