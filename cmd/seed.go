package cmd

import (
	"log"

	"github.com/h3poteto/fascia/db/seed"
	"github.com/spf13/cobra"
)

func seedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Insert seed data",
		Run:   seeds,
	}
	return cmd
}

func seeds(cmd *cobra.Command, args []string) {
	if err := seed.Seeds(); err != nil {
		log.Fatal(err)
	}
}
