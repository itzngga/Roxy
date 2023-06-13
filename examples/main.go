package main

import (
	_ "github.com/itzngga/roxy/examples/cmd"
	"log"

	"github.com/itzngga/roxy/core"
	"github.com/itzngga/roxy/options"
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
