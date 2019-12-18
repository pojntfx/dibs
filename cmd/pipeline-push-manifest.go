package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelinePushManifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Push manifest",
	Run: func(cmd *cobra.Command, args []string) {
		if output, err := Dibs.PushDockerManifest(viper.GetString(PlatformKey)); err != nil {
			utils.LogErrorFatal("Could not push manifest", err, output)
		}
	},
}

func init() {
	PipelinePushCmd.AddCommand(PipelinePushManifestCmd)
}
