package command

import (
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
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

func ButtonRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
	id, _ := args.Locals.Load("uid")
	button := util.CreateEmptyButton("Button", "@button",
		util.GenerateButton(id, "!help", "!help"),
		util.GenerateButton(id, "!hi", "!hi"),
		util.GenerateButton(id, "!sticker", "!sticker"))
	return button
}
