package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/alitto/pond"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/options"
	"github.com/itzngga/Roxy/types"
	"github.com/itzngga/Roxy/util"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waBinary "go.mau.fi/whatsmeow/binary"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"os"
	"time"
)

type App struct {
	Options *options.Options

	Log  waLog.Logger
	Pool *pond.WorkerPool

	startTime time.Time
	client    *whatsmeow.Client
	muxer     *Muxer
}

func NewGoRoxyBase(options options.Options) (*App, error) {
	err := options.Validate()
	if err != nil {
		return nil, err
	}

	optPointer := &options
	stdLog := waLog.Stdout("WaBOT", options.LogLevel, true)
	app := &App{
		Log:     stdLog,
		Options: optPointer,
		muxer:   NewMuxer(stdLog, optPointer),
	}
	app.Pool = pond.New(100, 1000)
	err = app.InitializeClient()
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (app *App) HandleEvents(event interface{}) {
	switch v := event.(type) {
	case *events.LoggedOut:
		app.Log.Warnf("%s Client logged out", app.client.Store.ID)
	case *events.Connected:
		app.startTime = time.Now()
		if len(app.client.Store.PushName) == 0 {
			return
		}
		app.Log.Infof("Connected!")
		_ = app.client.SendPresence(waTypes.PresenceAvailable)
	case *events.Message:
		app.Pool.Submit(func() {
			if !app.startTime.IsZero() && v.Info.Timestamp.After(app.startTime) {
				app.muxer.RunCommand(app.client, v)
				return
			}
		})
	case *events.StreamError:
		var message string
		if v.Code != "" {
			message = fmt.Sprintf("Unknown stream error with code %s", v.Code)
		} else if children := v.Raw.GetChildren(); len(children) > 0 {
			message = fmt.Sprintf("Unknown stream error (contains %s node)", children[0].Tag)
		} else {
			message = "Unknown stream error"
		}
		app.Log.Errorf("error: %s", message)
	case *events.CallOffer, *events.CallTerminate, *events.CallRelayLatency, *events.CallAccept, *events.UnknownCallEvent:
		// ignore
	case *events.AppState:
		// Ignore
	case *events.PushNameSetting:
		err := app.client.SendPresence(waTypes.PresenceAvailable)
		if err != nil {
			app.Log.Warnf("Failed to send presence after push name update: %v\n", err)
		}
	}
}
func (app *App) InitializeContainer() (*sqlstore.Container, error) {
	store.DeviceProps.RequireFullSync = types.Bool(false)
	if app.Options.WithSqlDB != nil {
		container := sqlstore.NewWithDB(app.Options.WithSqlDB, app.Options.StoreMode, waLog.Stdout("Database", "ERROR", true))
		err := container.Upgrade()
		if err != nil {
			panic(err)
		}
		return container, nil
	} else if app.Options.StoreMode == "postgres" {
		container, err := sqlstore.New("postgres", app.Options.PostgresDsn.GenerateDSN(), waLog.Stdout("Database", "ERROR", true))
		if err != nil {
			panic(err)
		}
		return container, nil
	} else if app.Options.StoreMode == "sqlite" {
		container, err := sqlstore.New("sqlite3", app.Options.SqliteFile, waLog.Stdout("Database", "ERROR", true))
		if err != nil {
			panic(err)
		}
		return container, nil
	} else {
		return nil, errors.New("error: invalid store mode")
	}
}

func (app *App) HandlePanic(p interface{}) {
	if p != nil {
		app.Log.Errorf("panic: \n%v", p)
	}
}

func (app *App) InitializeClient() error {
	container, err := app.InitializeContainer()
	if err != nil {
		return err
	}

	var device *store.Device
	if app.Options.HostNumber != "" {
		jid, ok := util.ParseJID(app.Options.HostNumber)
		if !ok {
			panic("invalid given number")
		}
		device, err = container.GetDevice(jid.ToNonAD())
		if err != nil {
			app.Log.Errorf("get device error %v", err)
		}
		if device == nil {
			device = container.NewDevice()
		}
	} else {
		device, err = container.GetFirstDevice()
		if err != nil {
			app.Log.Errorf("get device error %v", err)
		}
	}

	waBinary.IndentXML = true
	store.SetOSInfo("GoRoxy", [3]uint32{2, 2318, 11})
	store.DeviceProps.PlatformType = waProto.DeviceProps_CHROME.Enum()
	store.DeviceProps.RequireFullSync = types.Bool(false)

	app.client = whatsmeow.NewClient(device, waLog.Stdout("WhatsMeow", "ERROR", true))
	qrChan, _ := app.client.GetQRChannel(context.Background())
	err = app.client.Connect()
	if err != nil {
		if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			app.Log.Errorf("error connecting to client: %v", err)
		}
	} else {
		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					app.Log.Infof("QR Generated!")
				} else if evt.Event == "success" {
					app.Log.Infof("QR Scanned!")
				} else {
					app.Log.Infof("QR channel result: %s", evt.Event)
				}
			}
		}()
	}

	err = app.client.Connect()
	if err != nil {
		if !errors.Is(err, whatsmeow.ErrAlreadyConnected) {
			app.Log.Errorf("error connecting to client: %v", err)
			return err
		}
	}

	app.client.AddEventHandler(app.HandleEvents)
	return nil
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
