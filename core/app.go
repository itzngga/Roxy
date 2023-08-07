package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alitto/pond"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/options"
	"github.com/itzngga/Roxy/types"
	"github.com/itzngga/Roxy/util"
	"github.com/mdp/qrterminal/v3"
	"github.com/zhangyunhao116/skipmap"
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
	log     waLog.Logger
	pool    *pond.WorkerPool
	options *options.Options

	muxer        *Muxer
	sqlDB        *sql.DB
	pairCodeChan chan bool
	startTime    time.Time
	client       *whatsmeow.Client
	ctx          *skipmap.StringMap[types.RoxyContext]
}

func NewGoRoxyBase(options *options.Options) (*App, error) {
	stdLog := waLog.Stdout("WaBOT", options.LogLevel, true)
	app := &App{
		log:          stdLog,
		options:      options,
		pairCodeChan: make(chan bool),
	}
	app.pool = pond.New(100, 1000)
	err := app.initializeClient()
	if err != nil {
		return nil, err
	}
	if app.client != nil {
		app.generateContext()
		app.muxer = NewMuxer(app.ctx, stdLog, options)
	}

	return app, nil
}

func (app *App) handleEvents(event interface{}) {
	switch v := event.(type) {
	case *events.LoggedOut:
		app.log.Warnf("%s client logged out", app.client.Store.ID)
		app.client.Store.Delete()

		newApp, err := NewGoRoxyBase(app.options)
		if err != nil {
			panic(err)
		}
		*app = *newApp
	case *events.Connected:
		app.startTime = time.Now()
		if len(app.client.Store.PushName) == 0 {
			return
		}
		if app.options.LoginOptions == options.PAIR_CODE {
			app.pairCodeChan <- true
		}
		app.log.Infof("Connected!")
		app.ctx.Store("appClient", app.client)
		_ = app.client.SendPresence(waTypes.PresenceAvailable)
		app.muxer.cacheAllGroup()
	case *events.Message:
		app.pool.Submit(func() {
			if !app.startTime.IsZero() && v.Info.Timestamp.After(app.startTime) {
				app.muxer.RunCommand(app.client, v)
				return
			}
		})
		if app.options.HistorySync {
			app.pool.Submit(func() {
				app.handleMessageUpdates(v)
				return
			})
		}
	case *events.StreamError:
		var message string
		if v.Code != "" {
			message = fmt.Sprintf("Unknown stream error with code %s", v.Code)
		} else if children := v.Raw.GetChildren(); len(children) > 0 {
			message = fmt.Sprintf("Unknown stream error (contains %s node)", children[0].Tag)
		} else {
			message = "Unknown stream error"
		}
		app.log.Errorf("error: %s", message)
	case *events.CallOffer, *events.CallOfferNotice:
		var (
			callId string
			caller string
		)

		if val, ok := v.(*events.CallOffer); ok {
			callId = val.CallID
			caller = val.CallCreator.ToNonAD().String()
		} else if val, ok := v.(*events.CallOfferNotice); ok {
			callId = val.CallID
			caller = val.From.ToNonAD().String()
		}

		if app.options.AutoRejectCall {
			err := app.client.DangerousInternals().SendNode(waBinary.Node{
				Tag: "call",
				Attrs: waBinary.Attrs{
					"id":   whatsmeow.GenerateMessageID(),
					"from": app.client.Store.ID.ToNonAD().String(),
					"to":   caller,
				},
				Content: []waBinary.Node{
					{
						Tag: "reject",
						Attrs: waBinary.Attrs{
							"call-id":      callId,
							"call-creator": caller,
							"count":        "0",
						},
						Content: nil,
					},
				},
			})
			if err != nil {
				app.log.Errorf("failed to reject call: %v\n", err)
				return
			}
		}
	//case *events.CallOfferNotice:

	case *events.CallTerminate, *events.CallRelayLatency, *events.CallAccept, *events.UnknownCallEvent:
		// ignore
	case *events.AppState:
		// Ignore
	case *events.PushNameSetting:
		err := app.client.SendPresence(waTypes.PresenceAvailable)
		if err != nil {
			app.log.Warnf("Failed to send presence after push name update: %v\n", err)
		}
	case *events.JoinedGroup:
		app.muxer.unCacheOneGroup(nil, v)
	case *events.GroupInfo:
		app.muxer.unCacheOneGroup(v, nil)
	case *events.HistorySync:
		if app.options.HistorySync {
			app.pool.Submit(func() {
				app.handleHistorySync(v.Data)
				return
			})
		}
	}
}
func (app *App) initializeContainer() (*sqlstore.Container, error) {
	store.DeviceProps.RequireFullSync = types.Bool(app.options.HistorySync)
	if app.options.WithSqlDB != nil {
		container := sqlstore.NewWithDB(app.options.WithSqlDB, app.options.StoreMode, waLog.Stdout("Database", "ERROR", true))
		err := container.Upgrade()
		if err != nil {
			panic(err)
		}
		app.sqlDB = app.options.WithSqlDB
		app.initializeTables()
		return container, nil
	} else if app.options.StoreMode == "postgres" {
		db, err := sql.Open("postgres", app.options.PostgresDsn.GenerateDSN())
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}
		container := sqlstore.NewWithDB(db, "postgres", waLog.Stdout("Database", "ERROR", true))
		err = container.Upgrade()
		if err != nil {
			panic(err)
		}
		app.sqlDB = db
		app.initializeTables()
		return container, nil
	} else if app.options.StoreMode == "sqlite" {
		db, err := sql.Open("sqlite3", app.options.SqliteFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}
		container := sqlstore.NewWithDB(db, "sqlite3", waLog.Stdout("Database", "ERROR", true))
		err = container.Upgrade()
		if err != nil {
			panic(err)
		}
		app.sqlDB = db
		app.initializeTables()
		return container, nil
	} else {
		return nil, errors.New("error: invalid store mode")
	}
}

func (app *App) handlePanic(p interface{}) {
	if p != nil {
		app.log.Errorf("panic: \n%v", p)
	}
}

func (app *App) initializeClient() error {
	container, err := app.initializeContainer()
	if err != nil {
		return err
	}

	var device *store.Device
	if app.options.HostNumber != "" {
		jid, ok := util.ParseJID(app.options.HostNumber)
		if !ok {
			panic("invalid given number")
		}
		device, err = container.GetDevice(jid.ToNonAD())
		if err != nil {
			app.log.Errorf("get device error %v", err)
		}
		if device == nil {
			device = container.NewDevice()
		}
	} else {
		device, err = container.GetFirstDevice()
		if err != nil {
			app.log.Errorf("get device error %v", err)
		}
	}

	waBinary.IndentXML = true
	store.SetOSInfo("GoRoxy", [3]uint32{2, 2318, 11})
	store.DeviceProps.PlatformType = waProto.DeviceProps_CHROME.Enum()
	store.DeviceProps.RequireFullSync = types.Bool(false)

	app.client = whatsmeow.NewClient(device, waLog.Stdout("WhatsMeow", "ERROR", true))
	app.client.AddEventHandler(app.handleEvents)
	if app.options.LoginOptions == options.SCAN_QR {
		qrChan, _ := app.client.GetQRChannel(context.Background())
		err = app.client.Connect()
		if err != nil {
			if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
				app.log.Errorf("error connecting to client: %v", err)
			}
		} else {
			app.pool.Submit(func() {
				for evt := range qrChan {
					if evt.Event == "code" {
						qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
						app.log.Infof("QR Generated!")
					} else if evt.Event == "success" {
						app.log.Infof("QR Scanned!")
						app.ctx.Store("appClient", app.client)
					} else {
						app.log.Infof("QR channel result: %s", evt.Event)
					}
				}
			})
		}

		err = app.client.Connect()
		if err != nil {
			if !errors.Is(err, whatsmeow.ErrAlreadyConnected) {
				app.log.Errorf("error connecting to client: %v", err)
				return err
			}
		}
	} else {
		app.client.Disconnect()
		err = app.client.Connect()
		if err != nil {
			if !errors.Is(err, whatsmeow.ErrAlreadyConnected) {
				app.log.Errorf("error connecting to client: %v", err)
				return err
			}
		}
		if app.options.HostNumber == "" {
			return errors.New("error: you must specify host number when using pair code login options")
		}

		pairCode, err := app.client.PairPhone(app.options.HostNumber, true)
		if err != nil {
			return err
		}

		app.log.Infof("PairCode: %s", pairCode)

		for res := range app.pairCodeChan {
			if res {
				break
			}
		}
	}

	return nil
}

func (app *App) generateContext() {
	app.ctx = skipmap.NewString[types.RoxyContext]()
	app.ctx.Store("UpsertMessages", types.UpsertMessages(app.upsertMessages))
	app.ctx.Store("GetAllChats", types.GetAllChats(app.getAllChats))
	app.ctx.Store("GetChatInJID", types.GetChatInJID(app.getChatInJID))
	app.ctx.Store("GetStatusMessages", types.GetStatusMessages(app.getStatusMessages))
	app.ctx.Store("FindMessageByID", types.FindMessageByID(app.findMessageByID))
	app.ctx.Store("workerPool", app.pool)
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
