package command

import (
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func TestButtonCommand(c *whatsmeow.Client, m *events.Message, cmd *handler.Command) *waProto.Message {
	id := cmd.GetLocals("uid").(string)
	button := util.CreateButtonMessage("Button", "@button",
		util.GenerateButton(id, "!help", "!help"),
		util.GenerateButton(id, "!hi", "!hi"),
		util.GenerateButton(id, "!sticker", "!sticker"))
	return button
}
