package cmd

import "github.com/spf13/cobra"

// moduleCmd ist the entry point for the lang subcommand
var langCmd = &cobra.Command{
	Use:   "lang",
	Short: "Utilities for languages",
}

// langCmd is a subcommand of rootCmd
func init() {
	rootCmd.AddCommand(langCmd)
}
