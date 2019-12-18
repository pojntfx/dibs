package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PipelineTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test artifacts",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ReadConfig(viper.GetString(DibsFileKey))
	},
}

func init() {
	PipelineCmd.AddCommand(PipelineTestCmd)
}
