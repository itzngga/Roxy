package cmd

import (
	"fmt"
	waProto "github.com/go-whatsapp/whatsmeow/binary/proto"
	"github.com/itzngga/Roxy"
	"github.com/itzngga/Roxy/context"
	"time"
)

var speed = &roxy.Command{
	Name:        "speed",
	Description: "Testing speed",
	RunFunc: func(ctx *context.Ctx) *waProto.Message {
		t := time.Now()
		ctx.SendReplyMessage("wait...")
		return ctx.GenerateReplyMessage(fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
	},
}

func init() {
	roxy.Commands.Add(speed)
}
