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
	PrivateOnly: true,
	Category:    categories.CommonCategory,
	RunFunc: func(c *whatsmeow.Client, params *command.RunFuncParams) *waProto.Message {
		fromMe := false
		react := ""
		if len(params.Arguments) > 1 {
			react = params.Arguments[1]
		}
		if id := c.Store.ID.User + "@" + c.Store.ID.Server; *util.ParseQuotedRemoteJid(params.Event) == id {
			fromMe = true
		}
		return &waProto.Message{
			ReactionMessage: &waProto.ReactionMessage{
				Key: &waProto.MessageKey{
					RemoteJid: types.String(params.Event.Info.Chat.String()),
					FromMe:    &fromMe,
					Id:        util.ParseQuotedMessageId(params.Event),
				},
				Text:              types.String(react),
				SenderTimestampMs: types.Int64(time.Now().UnixMilli()),
			},
		}
	},
}

func init() {
	embed.Commands.Add(react)
}
