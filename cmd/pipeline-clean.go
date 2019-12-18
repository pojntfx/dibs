package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the output of pipeline building blocks",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ReadConfig(viper.GetString(DibsFileKey))
	},
}

func init() {
	PipelineCmd.AddCommand(PipelineCleanCmd)
}
