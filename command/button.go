package command

import (
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func ButtonCommand() {
	AddCommand(
		&handler.Command{
			Name:        "button",
			Category:    handler.UtilitiesCategory,
			Description: "Create button command",
			RunFunc:     ButtonRunFunc,
		})
}

func ButtonRunFunc(c *whatsmeow.Client, m *events.Message, cmd *handler.Command) *waProto.Message {
	id := cmd.GetLocals("uid").(string)
	button := util.CreateEmptyButton("Button", "@button",
		util.GenerateButton(id, "!help", "!help"),
		util.GenerateButton(id, "!hi", "!hi"),
		util.GenerateButton(id, "!sticker", "!sticker"))
	return button
}
