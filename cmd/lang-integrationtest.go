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
		if Platform == PlatformAll {
			switch On {
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
			switch On {
			case OnNative:
				if err := buildConfigs.TestIntegrationLang(Platform); err != nil {
					log.Error("Language integration test failed", rz.Err(err))
				}
			case OnDocker:
				if err := buildConfigs.TestIntegrationLangInDocker(Platform); err != nil {
					log.Error("Language integration test in Docker failed", rz.Err(err))
				}
			}
		}
	}}

func init() {
	langCmd.AddCommand(langIntegrationtestCmd)
}
