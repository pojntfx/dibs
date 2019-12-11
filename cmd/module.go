package cmd

import "github.com/spf13/cobra"

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Utilities for module development",
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}
