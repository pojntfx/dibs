package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var (
	NESTED bool
)

var binaryIntegrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the binary integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		if NESTED {
			if PLATFORM == PlatformAll {
				switch ON {
				case OnNative:
					if err := buildConfigs.TestIntegrationImageAll(); err != nil {
						log.Error("Nested integration tests failed", rz.Err(err))
					}
				case OnDocker:
					if err := buildConfigs.TestIntegrationImageInDockerAll(); err != nil {
						log.Error("Nested integration tests in Docker failed", rz.Err(err))
					}
				}
			} else {
				switch ON {
				case OnNative:
					if err := buildConfigs.TestIntegrationImage(PLATFORM); err != nil {
						log.Error("Nested integration test failed", rz.Err(err))
					}
				case OnDocker:
					if err := buildConfigs.TestIntegrationImageInDocker(PLATFORM); err != nil {
						log.Error("Nested integration tests in Docker failed", rz.Err(err))
					}
				}
			}
		} else {
			if PLATFORM == PlatformAll {
				switch ON {
				case OnNative:
					if err := buildConfigs.TestIntegrationBinaryAll(); err != nil {
						log.Error("Integration tests failed", rz.Err(err))
					}
				case OnDocker:
					if err := buildConfigs.TestIntegrationBinaryInDockerAll(); err != nil {
						log.Error("Integration tests in Docker failed", rz.Err(err))
					}
				}
			} else {
				switch ON {
				case OnNative:
					if err := buildConfigs.TestIntegrationBinary(PLATFORM); err != nil {
						log.Error("Integration test failed", rz.Err(err))
					}
				case OnDocker:
					if err := buildConfigs.TestIntegrationBinaryInDocker(PLATFORM); err != nil {
						log.Error("Integration test in Docker failed", rz.Err(err))
					}
				}
			}
		}
	}}

func init() {
	binaryIntegrationtestCmd.PersistentFlags().BoolVar(&NESTED, "nested", false, "Whether to integration test the image instead of the binary")

	binaryCmd.AddCommand(binaryIntegrationtestCmd)
}
