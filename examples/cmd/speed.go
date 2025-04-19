package cmd

import (
	"fmt"
	"github.com/itzngga/Roxy"
	"github.com/itzngga/Roxy/context"
	waProto "go.mau.fi/whatsmeow/binary/proto"
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
