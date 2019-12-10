package cmd

import "github.com/spf13/cobra"

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Utilities for Docker images",
}

func init() {
	rootCmd.AddCommand(imageCmd)
}
