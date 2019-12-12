package v2

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineBuildManifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Build manifest",
	Run: func(cmd *cobra.Command, args []string) {
		if output, err := Dibs.BuildDockerManifest(); err != nil {
			utils.PipeLogErrorNonPlatformSpecific("Could not build manifest", err, output)
		}
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildManifestCmd)
}
