package main

import (
	"fmt"
	"os"

	"github.com/h3poteto/fascia/cmd"
)

//go:generate go-assets-builder --output=config/bindata.go -s="/config" -p=config config/settings.yml

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
