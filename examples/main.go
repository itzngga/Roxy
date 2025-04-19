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
	opt.HostNumber = "6281395685501"
	opt.LoginOptions = options.PAIR_CODE

	app, err := roxy.NewRoxyBase(opt)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}
