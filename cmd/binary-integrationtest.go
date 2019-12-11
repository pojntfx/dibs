package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var (
	Image bool
)

var binaryIntegrationtestCmd = &cobra.Command{
	Use:   "integrationtest",
	Short: "Run the Docker image or binary integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		if Image {
			if Platform == PlatformAll {
				switch On {
				case OnNative:
					if err := buildConfigs.TestIntegrationImageAll(); err != nil {
						log.Error("Image integration tests failed", rz.Err(err))
					}
				case OnDocker:
					if err := buildConfigs.TestIntegrationImageInDockerAll(); err != nil {
						log.Error("Image integration tests in Docker failed", rz.Err(err))
					}
				}
			} else {
				switch On {
				case OnNative:
					if err := buildConfigs.TestIntegrationImage(Platform); err != nil {
						log.Error("Image integration test failed", rz.Err(err))
					}
				case OnDocker:
					if err := buildConfigs.TestIntegrationImageInDocker(Platform); err != nil {
						log.Error("Image integration tests in Docker failed", rz.Err(err))
					}
				}
			}
		} else {
			if Platform == PlatformAll {
				switch On {
				case OnNative:
					if err := buildConfigs.TestIntegrationBinaryAll(); err != nil {
						log.Error("Binary integration tests failed", rz.Err(err))
					}
				case OnDocker:
					if err := buildConfigs.TestIntegrationBinaryInDockerAll(); err != nil {
						log.Error("Integration tests in Docker failed", rz.Err(err))
					}
				}
			} else {
				switch On {
				case OnNative:
					if err := buildConfigs.TestIntegrationBinary(Platform); err != nil {
						log.Error("Integration test failed", rz.Err(err))
					}
				case OnDocker:
					if err := buildConfigs.TestIntegrationBinaryInDocker(Platform); err != nil {
						log.Error("Integration test in Docker failed", rz.Err(err))
					}
				}
			}
		}
	}}

func init() {
	binaryIntegrationtestCmd.PersistentFlags().BoolVar(&Image, "image", false, "Whether to work on images or binaries")

	binaryCmd.AddCommand(binaryIntegrationtestCmd)
}
