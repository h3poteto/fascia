package main

import (
	"github.com/h3poteto/fascia/server"
)

// TODO: subcomands

//go:generate go-bindata -ignore=\\.go -o=config/bindata.go -pkg=config -prefix=config/ config/

func main() {
	server.Serve()
}
