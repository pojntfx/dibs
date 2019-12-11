package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var langIntegrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		if PLATFORM == PlatformAll {
			switch ON {
			case OnNative:
				if err := buildConfigs.TestIntegrationLangAll(); err != nil {
					log.Error("Language integration tests failed", rz.Err(err))
				}
			case OnDocker:
				if err := buildConfigs.TestIntegrationLangInDockerAll(); err != nil {
					log.Error("Language integration tests in Docker failed", rz.Err(err))
				}
			}
		} else {
			switch ON {
			case OnNative:
				if err := buildConfigs.TestIntegrationLang(PLATFORM); err != nil {
					log.Error("Language integration test failed", rz.Err(err))
				}
			case OnDocker:
				if err := buildConfigs.TestIntegrationLangInDocker(PLATFORM); err != nil {
					log.Error("Language integration test in Docker failed", rz.Err(err))
				}
			}
		}
	}}

func init() {
	langCmd.AddCommand(langIntegrationtestCmd)
}
