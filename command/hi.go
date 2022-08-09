package command

import (
	"fmt"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

func HiCommand() {
	AddCommand(&handler.Command{
		Name:        "tes",
		Aliases:     []string{"hi"},
		Description: "A Fucking Test",
		Category:    handler.MiscCategory,
		RunFunc:     HiRunFunc,
	})
}

func HiRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
	t := time.Now()
	util.SendReplyMessage(c, args.Evm, "testing a...")
	return util.SendReplyText(args.Evm, fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
}
