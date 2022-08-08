package internal

import (
	"context"
	"fmt"
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type App struct {
	Log      *zap.Logger
	SqlStore *sqlstore.Container
}

type Base struct {
	startTime time.Time
	Log       *zap.Logger
	Muxer     *handler.Muxer
	client    *whatsmeow.Client
	Device    *store.Device
}

func (b *Base) QRChanFunc(ch <-chan whatsmeow.QRChannelItem) {
	for evt := range ch {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		} else {
			b.Log.Info("QRChan", zap.Any("QR channel result", evt.Event))
		}
	}
}

func (b *Base) Events(evt interface{}) {
	switch event := evt.(type) {
	case *events.Connected:
		b.startTime = time.Now()
		fmt.Println("[INFO] Connected!")
	case *events.Message:
		if !b.startTime.IsZero() && event.Info.Timestamp.After(b.startTime) {
			go b.Muxer.RunCommand(b.client, event)
		}
	}
}

func (b *Base) Init() {
	store.DeviceProps.RequireFullSync = proto.Bool(true)
	b.Muxer = handler.NewMuxer()

	for _, cmd := range command.Commands {
		b.Muxer.AddCommand(cmd)
	}

	b.client = whatsmeow.NewClient(b.Device, config.NewWaLog())
	if b.client.Store.ID == nil {
		qrChan, _ := b.client.GetQRChannel(context.Background())
		err := b.client.Connect()
		if err != nil {
			b.Log.With(zap.Error(err)).Error(err.Error())
		}
		go b.QRChanFunc(qrChan)
	} else {
		err := b.client.Connect()
		if err != nil {
			b.Log.With(zap.Error(err)).Error(err.Error())
		}
	}
	b.client.AddEventHandler(b.Events)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	b.client.Disconnect()
}
