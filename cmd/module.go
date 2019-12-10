package cmd

import "github.com/spf13/cobra"

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Utilities for modules",
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}
