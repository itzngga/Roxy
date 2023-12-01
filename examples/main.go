package main

import (
	roxy "github.com/itzngga/Roxy"
	_ "github.com/itzngga/Roxy/examples/cmd"
	"github.com/itzngga/Roxy/options"

	_ "github.com/mattn/go-sqlite3"

	"log"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	opt := options.NewDefaultOptions()
	app, err := roxy.NewRoxyBase(opt)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}
