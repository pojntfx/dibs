package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineBuildManifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Build manifest",
	Run: func(cmd *cobra.Command, args []string) {
		if output, err := Dibs.BuildDockerManifest(Platform); err != nil {
			utils.PipeLogErrorFatalNonPlatformSpecific("Could not build manifest", err, output)
		}
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildManifestCmd)
}
