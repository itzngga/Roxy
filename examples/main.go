package main

import (
	roxy "github.com/itzngga/Roxy"
	_ "github.com/itzngga/Roxy/examples/cmd"
	"github.com/itzngga/Roxy/options"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	opt := options.NewDefaultOptions()
	opt.HostNumber = os.Getenv("HOST_NUMBER")
	opt.LoginOptions = options.PAIR_CODE
	opt.HistorySync = true

	app, err := roxy.NewRoxyBase(opt)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}
