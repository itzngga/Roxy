package cmd

import (
	"fmt"
	"github.com/itzngga/goRoxy/basic/categories"
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/embed"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

var speed = &command.Command{
	Name:        "speed",
	Aliases:     []string{"sp", "s"},
	Description: "Testing speed",
	Category:    categories.CommonCategory,
	RunFunc: func(c *whatsmeow.Client, args command.RunFuncArgs) *waProto.Message {
		t := time.Now()
		util.SendReplyMessage(c, args.Evm, "ok, waitt...")
		return util.SendReplyText(args.Evm, fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
	},
}

func init() {
	embed.Commands.Add(speed)
}
