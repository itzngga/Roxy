package main

import (
	_ "github.com/itzngga/Roxy/examples/cmd"
	"log"

	"github.com/itzngga/Roxy/core"
	"github.com/itzngga/Roxy/options"
	_ "github.com/mattn/go-sqlite3"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := core.NewGoRoxyBase(options.NewDefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}
