package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var binaryBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the binary",
	Run: func(cmd *cobra.Command, args []string) {
		if PLATFORM == PlatformAll {
			switch ON {
			case OnNative:
				if err := buildConfigs.BuildAll(); err != nil {
					log.Error("Could not build the binaries", rz.Err(err))
				}
			case OnDocker:
				if err := buildConfigs.BuildInDockerAll(); err != nil {
					log.Error("Could not build the binaries in Docker", rz.Err(err))
				}
				if err := buildConfigs.GetBinaryFromDockerImageAll(); err != nil {
					log.Error("Could not get binaries from Docker images", rz.Err(err))
				}
			}
		} else {
			switch ON {
			case OnNative:
				if err := buildConfigs.Build(PLATFORM); err != nil {
					log.Error("Could not build the binary", rz.Err(err))
				}
			case OnDocker:
				if err := buildConfigs.BuildInDocker(PLATFORM); err != nil {
					log.Error("Could not build the binary in Docker", rz.Err(err))
				}
				if err := buildConfigs.GetBinaryFromDockerImageAll(); err != nil {
					log.Error("Could not get binary from Docker image", rz.Err(err))
				}
			}
		}
	}}

func init() {
	binaryCmd.AddCommand(binaryBuildCmd)
}
