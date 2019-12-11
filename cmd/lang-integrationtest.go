package cmd

import "github.com/spf13/cobra"

var langIntegrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		if PLATFORM == PLATFORM_ALL {
			switch ON {
			case ON_NATIVE:
				buildConfigs.TestIntegrationLangAll()
			case ON_DOCKER:
				buildConfigs.TestIntegrationLangInDockerAll()
			}
		} else {
			switch ON {
			case ON_NATIVE:
				buildConfigs.TestIntegrationLang(PLATFORM)
			case ON_DOCKER:
				buildConfigs.TestIntegrationLangInDocker(PLATFORM)
			}
		}
	},
}

func init() {
	langCmd.AddCommand(langIntegrationtestCmd)
}
