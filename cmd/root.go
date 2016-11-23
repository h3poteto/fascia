package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd has commands for fascia
var RootCmd = &cobra.Command{
	Use:           "fascia",
	Short:         "Fascia server commands",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	cobra.OnInitialize()
	RootCmd.AddCommand(
		serverCmd(),
		seedCmd(),
	)
}
