package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var binaryCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the Docker image or binary output",
	Run: func(cmd *cobra.Command, args []string) {
		if Platform == PlatformAll {
			switch On {
			case OnNative:
				if err := buildConfigs.BuildCleanAll(); err != nil {
					log.Error("Could not clean binaries", rz.Err(err))
				}
			case OnDocker:
				if err := buildConfigs.BuildInDockerCleanAll(); err != nil {
					log.Error("Could not clean binaries in Docker", rz.Err(err))
				}
			}
		} else {
			switch On {
			case OnNative:
				if err := buildConfigs.BuildClean(Platform); err != nil {
					log.Error("Could not clean binary", rz.Err(err))
				}
			case OnDocker:
				if err := buildConfigs.BuildInDockerClean(Platform); err != nil {
					log.Error("Could not clean binary in Docker", rz.Err(err))
				}
			}
		}
	},
}

func init() {
	binaryCmd.AddCommand(binaryCleanCmd)
}
