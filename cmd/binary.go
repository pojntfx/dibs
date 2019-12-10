package cmd

import "github.com/spf13/cobra"

var binaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "Utilities for binaries",
}

func init() {
	rootCmd.AddCommand(binaryCmd)
}
