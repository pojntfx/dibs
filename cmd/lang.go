package cmd

import "github.com/spf13/cobra"

var langCmd = &cobra.Command{
	Use:   "lang",
	Short: "Utilities for languages",
}

func init() {
	rootCmd.AddCommand(langCmd)
}
