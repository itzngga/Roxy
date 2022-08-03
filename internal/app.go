package internal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itzngga/goRoxy/config"
	"github.com/itzngga/goRoxy/helper"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type AppConfig struct {
	waLog    waLog.Logger
	log      *zap.Logger
	sqlstore *sqlstore.Container
}

type Base struct {
	startTime time.Time
	client    *whatsmeow.Client
	device    *store.Device
	appConfig *AppConfig
}

func (b *Base) QRChanFunc(ch <-chan whatsmeow.QRChannelItem) {
	for evt := range ch {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		} else {
			b.appConfig.waLog.Infof("QR channel result: %s", evt.Event)
		}
	}
}

func (b *Base) Events() {
	b.client.AddEventHandler(func(evt interface{}) {
		switch event := evt.(type) {
		case *events.Connected:
			b.startTime = time.Now()
		case *events.Message:
			if !b.startTime.IsZero() && time.Now().After(b.startTime) {
				go handler.RunCommand(b.client, event)
				return
			}
			return
		}
	})
}

func (b *Base) initConfig() *AppConfig {
	return &AppConfig{
		waLog:    config.NewWaLog(),
		log:      config.NewLogger("ingfo"),
		sqlstore: config.SqlstoreContainer(),
	}
}

func (b *Base) Init() {
	store.DeviceProps.RequireFullSync = proto.Bool(true)
	b.appConfig = b.initConfig()
	b.InitializeCommands()
	if helper.XTCrypto != nil {
		helper.XTCrypto = helper.NewCryptography()
	}

	device, err := b.appConfig.sqlstore.GetFirstDevice()
	if err != nil {
		b.appConfig.waLog.Errorf("Failed to get device: %v", err)
	}
	b.device = device
	b.client = whatsmeow.NewClient(b.device, b.appConfig.waLog)
	if b.client.Store.ID == nil {
		qrChan, _ := b.client.GetQRChannel(context.Background())
		err := b.client.Connect()
		if err != nil {
			panic(err)
		}
		go b.QRChanFunc(qrChan)
	} else {
		err = b.client.Connect()
		if err != nil {
			panic(err)
		}
	}
	b.Events()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	b.client.Disconnect()

}
