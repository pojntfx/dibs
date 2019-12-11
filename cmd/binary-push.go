package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var binaryPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push Docker images or binaries",
	Run: func(cmd *cobra.Command, args []string) {
		if Image {
			if Platform == PlatformAll {
				if err := buildConfigs.PushDockerImageAll(); err != nil {
					log.Error("Could not push Docker images", rz.Err(err))
				}
			} else {
				if err := buildConfigs.PushDockerImage(Platform); err != nil {
					log.Error("Could not push Docker image", rz.Err(err))
				}
			}
		} else {
			log.Fatal("Not yet implemented")
		}
	},
}

func init() {
	binaryPushCmd.PersistentFlags().BoolVar(&Image, "image", false, "Whether to work on images or binaries")

	binaryCmd.AddCommand(binaryPushCmd)
}
