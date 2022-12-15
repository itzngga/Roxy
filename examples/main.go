package main

import (
	_ "github.com/itzngga/goRoxy/examples/cmd"

	"github.com/itzngga/goRoxy/core"
	"github.com/itzngga/goRoxy/options"
	_ "github.com/lib/pq"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := core.NewGoRoxyBase(options.NewDefaultOptions())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}
