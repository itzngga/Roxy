/*
Copyright © 2022 itzngga rangganak094@gmail.com. All rights reserved
*/
package main

import (
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/config"
	"github.com/itzngga/goRoxy/internal"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"go.uber.org/zap"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	handler.GlobalLocals = &map[string]interface{}{}
	command.Commands = make([]*handler.Command, 0)
	handler.GlobalMiddleware = make([]handler.MiddlewareFunc, 0)
	command.GenerateAllCommands()
	middleware.GenerateAllMiddlewares()
}

func main() {
	app := &internal.App{
		Log:      config.NewLogger("info"),
		SqlStore: config.SqlStoreContainer(),
	}

	device, err := app.SqlStore.GetFirstDevice()
	if err != nil {
		app.Log.With(zap.Error(err)).Error(err.Error())
	}

	base := internal.Base{
		Device: device,
		Log:    app.Log,
	}

	base.Init()
}
