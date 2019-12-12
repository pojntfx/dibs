package cmd

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelinePushManifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Push manifest",
	Run: func(cmd *cobra.Command, args []string) {
		if output, err := Dibs.PushDockerManifest(Platform); err != nil {
			utils.PipeLogErrorFatalNonPlatformSpecific("Could not push manifest", err, output)
		}
	},
}

func init() {
	PipelinePushCmd.AddCommand(PipelinePushManifestCmd)
}
