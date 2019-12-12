package v2

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var PipelineBuildImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Build image",
	Run: func(cmd *cobra.Command, args []string) {
		for _, platform := range Dibs.Platforms {
			if Platform == PlatformAll {
				if output, err := platform.Binary.Build.BuildImage(platform.Platform); err != nil {
					log.Fatal("Could not build image", rz.String("platform", platform.Platform), rz.String("output", output), rz.Err(err))
				}
				return
			} else {
				if platform.Platform == Platform {
					if output, err := platform.Binary.Build.BuildImage(platform.Platform); err != nil {
						log.Fatal("Could not build image", rz.String("platform", platform.Platform), rz.String("output", output), rz.Err(err))
					}
					return
				}
			}
		}
		log.Fatal("Platform(s) not found in configuration file", rz.Any("platform", Platform))
	},
}

func init() {
	PipelineBuildCmd.AddCommand(PipelineBuildImageCmd)
}
