package cmd

import "github.com/spf13/cobra"

var imageCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the image output",
	Run: func(cmd *cobra.Command, args []string) {
		buildConfigs.Clean(PLATFORM)
	},
}

func init() {
	imageCmd.AddCommand(imageCleanCmd)
}
