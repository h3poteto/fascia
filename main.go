package main

import (
	"fmt"
	"os"

	"github.com/h3poteto/fascia/cmd"
)

//go:generate go-bindata -ignore=\\.go -o=config/bindata.go -pkg=config -prefix=config/ config/

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
