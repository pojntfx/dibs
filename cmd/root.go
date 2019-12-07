package cmd

import (
	"github.com/spf13/cobra"
	rz "gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

var rootCmd = &cobra.Command{
	Use:   "godibs",
	Short: "Distributed build system for Go",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error", rz.String("System", "Client"), rz.Err(err))
	}
}
