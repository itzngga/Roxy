package core

import (
	"context"
	"github.com/alitto/pond"
	_ "github.com/itzngga/goRoxy/basic/categories"
	_ "github.com/itzngga/goRoxy/basic/commands"
	_ "github.com/itzngga/goRoxy/basic/global_middleware"
	"github.com/itzngga/goRoxy/command"
	_ "github.com/itzngga/goRoxy/embed"
	"github.com/itzngga/goRoxy/options"
	"github.com/itzngga/goRoxy/types"
	"github.com/itzngga/goRoxy/util"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"os"
	"time"
)

type App struct {
	Options *options.Options

	SqlStore *sqlstore.Container
	Log      waLog.Logger
	Pool     *pond.WorkerPool

	startTime time.Time
	client    *whatsmeow.Client
	muxer     *Muxer
}

func NewGoRoxyBase(options options.Options) *App {
	options.Validate()
	optPointer := &options
	stdLog := waLog.Stdout("WaBOT", options.LogLevel, true)
	app := &App{
		Log:     stdLog,
		Options: optPointer,
		muxer:   NewMuxer(stdLog, optPointer),
	}
	app.Pool = pond.New(100, 1000)
	app.PrepareClient()
	return app
}

func (app *App) QRChanFunc(ch <-chan whatsmeow.QRChannelItem) {
	for evt := range ch {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		} else {
			app.Log.Infof("[INFO] QR channel result %a", evt.Event)
		}
	}
}

func (app *App) ConnectedEvents(evt interface{}) {
	_, ok := evt.(*events.Connected)
	if ok {
		app.startTime = time.Now()
		app.Log.Infof("Connected!")
	}
}

func (app *App) MessageEvents(evt interface{}) {
	event, ok := evt.(*events.Message)
	if ok {
		app.Pool.Submit(func() {
			if !app.startTime.IsZero() && event.Info.Timestamp.After(app.startTime) {
				app.muxer.RunCommand(app.client, event)
				return
			}
		})
	}
	return
}
func (app *App) PrepareSqlContainer() {
	store.DeviceProps.RequireFullSync = types.Bool(true)
	if app.Options.StoreMode == "postgres" {
		container, err := sqlstore.New("postgres", app.Options.PostgresDsn, waLog.Stdout("Database", "ERROR", true))
		if err != nil {
			panic(err)
		}
		app.SqlStore = container
	} else {
		container, err := sqlstore.New("sqlite3", app.Options.SqliteFile, waLog.Stdout("Database", "ERROR", true))
		if err != nil {
			panic(err)
		}
		app.SqlStore = container
	}
}

func (app *App) HandlePanic(p interface{}) {
	if p != nil {
		app.Log.Errorf("panic: \n%v", p)
	}
}

func (app *App) PrepareClient() {
	app.PrepareSqlContainer()

	var device *store.Device
	var err error
	if app.Options.HostNumber != "" {
		jid, ok := util.ParseJID(app.Options.HostNumber)
		if !ok {
			panic("invalid given number")
		}
		device, err = app.SqlStore.GetDevice(jid)
		if err != nil {
			app.Log.Errorf("get device error %v", err)
		}
	} else {
		device, err = app.SqlStore.GetFirstDevice()
		if err != nil {
			app.Log.Errorf("get device error %v", err)
		}
	}

	app.client = whatsmeow.NewClient(device, waLog.Stdout("WhatsMeow", "ERROR", true))
	if app.client.Store.ID == nil {
		qrChan, _ := app.client.GetQRChannel(context.Background())
		err := app.client.Connect()
		if err != nil {
			app.Log.Errorf("error connecting to client: %v", err)
		}
		go app.QRChanFunc(qrChan)
	} else {
		err := app.client.Connect()
		if err != nil {
			app.Log.Errorf("error connecting to client: %v", err)
		}
	}

	app.client.AddEventHandler(app.ConnectedEvents)
	app.client.AddEventHandler(app.MessageEvents)
}

func (app *App) AddNewCategory(category string) {
	app.muxer.Categories.Store(category, category)
}

func (app *App) AddNewCommand(command command.Command) {
	app.muxer.AddCommand(&command)
}

func (app *App) AddNewMiddleware(middleware command.MiddlewareFunc) {
	app.muxer.AddMiddleware(middleware)
}

func (app *App) Shutdown() {
	app.client.Disconnect()
}
