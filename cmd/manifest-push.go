package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var manifestPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push the Docker manifest",
	Run: func(cmd *cobra.Command, args []string) {
		if err := buildConfigs.PushDockerManifest(); err != nil {
			log.Error("Could not push the Docker manifest", rz.Err(err))
		}
	},
}

func init() {
	manifestCmd.AddCommand(manifestPushCmd)
}
