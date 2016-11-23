package cmd

import (
	"github.com/h3poteto/fascia/server"
	"github.com/spf13/cobra"
)

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start fascia server",
		Run:   serve,
	}
	return cmd
}

func serve(cmd *cobra.Command, args []string) {
	server.Serve()
}
