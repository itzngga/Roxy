package cmd

import (
	"fmt"
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/embed"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

var speed = &command.Command{
	Name:        "speed",
	Description: "Testing speed",
	//Category:    categories.CommonCategory,
	RunFunc: func(c *whatsmeow.Client, params *command.RunFuncParams) *waProto.Message {
		t := time.Now()
		util.SendReplyMessage(c, params.Event, "wait...")
		return util.GenerateReplyMessage(params.Event, fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
	},
}

func init() {
	embed.Commands.Add(speed)
}
