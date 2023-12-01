package main

import (
	roxy "github.com/itzngga/Roxy"
	_ "github.com/itzngga/Roxy/examples/cmd"
	"github.com/itzngga/Roxy/options"
	"log"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	pg := options.NewPostgresDSN()
	pg.SetHost("localhost")
	pg.SetPort("4321")
	pg.SetUsername("postgres")
	pg.SetPassword("root123")
	pg.SetTimeZone("Asia/Jakarta")

	opt := options.NewDefaultOptions()
	opt.StoreMode = "postgres"
	opt.PostgresDsn = pg

	app, err := roxy.NewRoxyBase(opt)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	app.Shutdown()
}
