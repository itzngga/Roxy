package roxy

import (
	contextCtx "context"
	"errors"
	"fmt"
	"github.com/go-whatsapp/whatsmeow"
	waBinary "github.com/go-whatsapp/whatsmeow/binary"
	waProto "github.com/go-whatsapp/whatsmeow/binary/proto"
	"github.com/go-whatsapp/whatsmeow/store"
	waTypes "github.com/go-whatsapp/whatsmeow/types"
	"github.com/go-whatsapp/whatsmeow/types/events"
	waLog "github.com/go-whatsapp/whatsmeow/util/log"
	"github.com/itzngga/Roxy/container"
	"github.com/itzngga/Roxy/context"
	"github.com/itzngga/Roxy/options"
	"github.com/mdp/qrterminal/v3"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	log       waLog.Logger
	options   *options.Options
	container *container.Container

	muxer        *Muxer
	startTime    time.Time
	device       *store.Device
	client       *whatsmeow.Client
	clientJID    waTypes.JID
	pairCodeChan chan bool
}

func NewRoxyBase(options *options.Options) (*App, error) {
	err := options.Validate()
	if err != nil {
		return nil, err
	}
	stdLog := waLog.Stdout("WaBOT", options.LogLevel, true)
	app := &App{
		log:          stdLog,
		options:      options,
		pairCodeChan: make(chan bool),
	}
	err = app.InitializeClient()
	if err != nil {
		return nil, err
	}
	if app.client != nil {
		app.muxer = NewMuxer(stdLog, options, app)
		app.muxer.addEmbedCommands()
	}

	return app, nil
}

func (app *App) InitializeClient() error {
	containerSql, err := container.NewContainer(app.options)
	if err != nil {
		return err
	}

	device, err := containerSql.AcquireDevice(app.options.HostNumber)
	if err != nil {
		return err
	}

	app.container = containerSql
	app.device = device

	waBinary.IndentXML = true
	store.SetOSInfo(app.options.OSInfo, store.GetWAVersion())
	store.DeviceProps.PlatformType = waProto.DeviceProps_DESKTOP.Enum()

	app.client = whatsmeow.NewClient(app.device, waLog.Stdout("WhatsMeow", "ERROR", true))
	app.client.EnableAutoReconnect = true
	app.client.AutoTrustIdentity = true
	app.client.AutomaticMessageRerequestFromPhone = true
	app.client.AddEventHandler(app.HandleEvents)

	if !app.container.NewDevice {
		err = app.client.Connect()
		return err
	}

	if app.options.LoginOptions == options.SCAN_QR {
		qrChan, _ := app.client.GetQRChannel(contextCtx.Background())
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				app.log.Infof("QR Generated!")
			} else if evt.Event == "success" {
				app.log.Infof("QR Scanned!")
			} else {
				app.log.Infof("QR channel result: %s", evt.Event)
			}
		}

		err = app.client.Connect()
		if err != nil {
			if !errors.Is(err, whatsmeow.ErrAlreadyConnected) {
				app.log.Errorf("error connecting to client: %v", err)
				return err
			}
		}
	} else {
		if app.options.HostNumber == "" {
			return errors.New("error: you must spcify host number when using pair code login options")
		}

		pairCode, err := app.client.PairPhone(app.options.HostNumber, true, whatsmeow.PairClientFirefox, "Firefox (Linux)")
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

func (app *App) HandleEvents(event interface{}) {
	switch v := event.(type) {
	case *events.LoggedOut:
		app.log.Warnf("%s client logged out", app.clientJID)
		app.client.Store.Delete()

		newApp, err := NewRoxyBase(app.options)
		if err != nil {
			panic(err)
		}
		*app = *newApp
	case *events.PairSuccess:
		app.clientJID = v.ID
	case *events.Connected:
		app.startTime = time.Now()
		if len(app.client.Store.PushName) == 0 {
			return
		}
		if app.options.LoginOptions == options.PAIR_CODE {
			app.pairCodeChan <- true
		}
		app.log.Infof("Connected!")
		app.clientJID = *app.client.Store.JID
		app.container.SetClientJID(app.clientJID)
		_ = app.client.SendPresence(waTypes.PresenceAvailable)
		//app.muxer.CacheAllGroup()
	case *events.Message:
		go func() {
			if !app.startTime.IsZero() && v.Info.Timestamp.After(app.startTime) {
				app.muxer.RunCommand(app.client, v)
				return
			}
		}()
		if app.options.HistorySync {
			app.container.HandleMessageUpdates(v)
			return
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
					"from": app.clientJID.ToNonAD().String(),
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
		app.muxer.UnCacheOneGroup(nil, v)
	case *events.GroupInfo:
		app.muxer.UnCacheOneGroup(v, nil)
	case *events.HistorySync:
		if app.options.HistorySync {
			app.container.HandleHistorySync(v.Data)
			return
		}
	}
}

func (app *App) handlePanic(p interface{}) {
	if p != nil {
		app.log.Errorf("panic: \n%v", p)
	}
}

func (app *App) Client() *whatsmeow.Client {
	return app.client
}

func (app *App) ClientJID() waTypes.JID {
	return app.clientJID
}

func (app *App) AddNewCategory(category string) {
	app.muxer.Categories.Store(category, category)
}

func (app *App) AddNewCommand(command Command) {
	app.muxer.AddCommand(&command)
}

func (app *App) AddNewMiddleware(middleware context.MiddlewareFunc) {
	app.muxer.AddMiddleware(middleware)
}

func (app *App) Shutdown() {
	container.ReleaseContainer(app.container)
	app.client.Disconnect()
	app.muxer = nil
}

func (app *App) SendMessage(to waTypes.JID, message *waProto.Message, extra ...whatsmeow.SendRequestExtra) (whatsmeow.SendResponse, error) {
	ctx, cancel := contextCtx.WithTimeout(contextCtx.Background(), app.options.SendMessageTimeout)
	defer cancel()

	response, err := app.client.SendMessage(ctx, to, message, extra...)
	if err != nil {
		app.log.Errorf("send message error: %v", err)
	}

	return response, err
}

func (app *App) UpsertMessages(jid waTypes.JID, message []*events.Message) {
	app.container.UpsertMessages(jid, message)
	return
}

func (app *App) GetAllChats() []*events.Message {
	return app.container.GetAllChats()
}

func (app *App) GetChatInJID(jid waTypes.JID) []*events.Message {
	return app.container.GetChatInJID(jid)
}

func (app *App) GetStatusMessages() []*events.Message {
	return app.container.GetStatusMessages()
}

func (app *App) FindMessageByID(jid waTypes.JID, id string) *events.Message {
	return app.container.FindMessageByID(jid, id)
}
