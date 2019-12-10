package cmd

import "github.com/spf13/cobra"

// moduleCmd ist the entry point for the module subcommand
var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Utilities for modules",
}

// moduleCmd is a subcommand of rootCmd
func init() {
	rootCmd.AddCommand(moduleCmd)
}
