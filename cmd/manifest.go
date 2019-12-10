package cmd

import "github.com/spf13/cobra"

var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Utilities for Docker manifests",
}

func init() {
	rootCmd.AddCommand(manifestCmd)
}
