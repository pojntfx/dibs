package v2

import (
	"github.com/pojntfx/godibs/pkg/utils"
	"github.com/spf13/cobra"
)

var PipelineBuildBinaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "Build binary",
	Run: func(cmd *cobra.Command, args []string) {
		platforms, err := Dibs.GetPlatforms(Platform, Platform == PlatformAll)
		if err != nil {
			utils.PipeLogErrorFatalPlatformNotFound(platforms, err)
		}

		for _, platform := range platforms {
			if Executor == ExecutorDocker {
				if output, err := platform.Binary.GetBinaryFromDockerImage(platform.Platform); err != nil {
					utils.PipeLogErrorFatal("Could not get binary from Docker image", err, platform.Platform, output)
				}
			} else {
				if output, err := platform.Binary.Build.Start(platform.Platform); err != nil {
					utils.PipeLogErrorFatal("Could not build binary", err, platform.Platform, output)
				}
			}
		}
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildBinaryCmd)
}
