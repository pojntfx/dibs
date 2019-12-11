package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var langCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the language output",
	Run: func(cmd *cobra.Command, args []string) {
		if err := buildConfigs.Clean(PLATFORM); err != nil {
			log.Error("Could not clean the language output", rz.Err(err))
		}
	},
}

func init() {
	langCmd.AddCommand(langCleanCmd)
}
