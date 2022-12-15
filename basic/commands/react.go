package commands

import (
	"github.com/itzngga/goRoxy/basic/categories"
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/embed"
	"github.com/itzngga/goRoxy/types"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

var react = &command.Command{
	Name:        "react",
	Description: "Emoji ke reply",
	Category:    categories.CommonCategory,
	RunFunc: func(c *whatsmeow.Client, args command.RunFuncArgs) *waProto.Message {
		react := ""
		if len(args.Args) != 1 {
			react = args.Args[1]
		}
		msg := &waProto.Message{
			ReactionMessage: &waProto.ReactionMessage{
				Key: &waProto.MessageKey{
					RemoteJid: types.String(args.Evm.Info.Chat.String()),
					FromMe:    types.Bool(false),
					Id:        util.ParseQuotedMessageId(args.Evm.Message),
				},
				Text:              types.String(react),
				SenderTimestampMs: types.Int64(time.Now().UnixMilli()),
			},
		}

		return msg
	},
}

func init() {
	embed.Commands.Add(react)
}
