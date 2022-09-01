/*
Copyright © 2022 itzngga rangganak094@gmail.com. All rights reserved
*/
package main

import (
	"fmt"
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/config"
	"github.com/itzngga/goRoxy/internal"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/middleware"
	"github.com/itzngga/goRoxy/util/cmd_gen"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zhangyunhao116/skipmap"
	"go.uber.org/zap"

	_ "modernc.org/sqlite"
)

func init() {
	cmd_gen.GenCmd()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	command.Commands = make([]*handler.Command, 0)
	handler.GlobalMiddleware = skipmap.NewString[handler.MiddlewareFunc]()
	command.GenerateAllCommands()
	middleware.GenerateAllMiddlewares()
}

func main() {
	fmt.Println("[INFO] Done generating commands")
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
