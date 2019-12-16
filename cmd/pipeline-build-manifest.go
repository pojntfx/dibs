package cmd

import (
	"github.com/pojntfx/dibs/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineBuildManifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Build manifest",
	Run: func(cmd *cobra.Command, args []string) {
		if output, err := Dibs.BuildDockerManifest(viper.GetString(PlatformKey)); err != nil {
			utils.PipeLogErrorFatalNonPlatformSpecific("Could not build manifest", err, output)
		}
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildManifestCmd)
}
