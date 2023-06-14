package cmd

import (
	"fmt"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

var speed = &command.Command{
	Name:        "speed",
	Description: "Testing speed",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		t := time.Now()
		ctx.SendReplyMessage("wait...")
		return ctx.GenerateReplyMessage(fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
	},
}

func init() {
	embed.Commands.Add(speed)
}
