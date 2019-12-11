package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var manifestBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Docker manifest",
	Run: func(cmd *cobra.Command, args []string) {
		if err := buildConfigs.BuildDockerManifest(); err != nil {
			log.Error("Could not build Docker manifest", rz.Err(err))
		}
	},
}

func init() {
	manifestCmd.AddCommand(manifestBuildCmd)
}
