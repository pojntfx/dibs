package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var binaryPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push the Docker image",
	Run: func(cmd *cobra.Command, args []string) {
		if NESTED {
			if PLATFORM == PlatformAll {
				if err := buildConfigs.PushDockerImageAll(); err != nil {
					log.Error("Could not push Docker images", rz.Err(err))
				}
			} else {
				if err := buildConfigs.PushDockerImage(PLATFORM); err != nil {
					log.Error("Could not push Docker image", rz.Err(err))
				}
			}
		} else {
			log.Fatal("Not yet implementated")
		}
	},
}

func init() {
	binaryPushCmd.PersistentFlags().BoolVar(&NESTED, "nested", false, "Whether to push the image instead of the binary")

	binaryCmd.AddCommand(binaryPushCmd)
}
