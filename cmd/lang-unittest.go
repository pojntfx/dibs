package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var langUnittestCmd = &cobra.Command{
	Use:   "unittest",
	Short: "Run the unit tests",
	Run: func(cmd *cobra.Command, args []string) {
		if Platform == PlatformAll {
			switch On {
			case OnNative:
				if err := buildConfigs.TestUnitAll(); err != nil {
					log.Error("Unit tests failed", rz.Err(err))
				}
			case OnDocker:
				if err := buildConfigs.TestUnitInDockerAll(); err != nil {
					log.Error("Unit tests failed in Docker", rz.Err(err))
				}
			}
		} else {
			switch On {
			case OnNative:
				if err := buildConfigs.TestUnit(Platform); err != nil {
					log.Error("Unit test failed", rz.Err(err))
				}
			case OnDocker:
				if err := buildConfigs.TestUnitInDocker(Platform); err != nil {
					log.Error("Unit test failed in Docker", rz.Err(err))
				}
			}
		}
	}}

func init() {
	langCmd.AddCommand(langUnittestCmd)
}
