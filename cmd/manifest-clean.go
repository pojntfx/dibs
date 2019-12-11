package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var manifestCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the Docker manifest output",
	Run: func(cmd *cobra.Command, args []string) {
		if err := buildConfigs.Clean(PLATFORM); err != nil {
			log.Error("Could not clean the Docker manifest output", rz.Err(err))
		}
	},
}

func init() {
	manifestCmd.AddCommand(manifestCleanCmd)
}
