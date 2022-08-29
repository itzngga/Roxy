package command

import (
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func HydratedCommand() {
	AddCommand(
		&handler.Command{
			Name:        "hydrated",
			Category:    handler.UtilitiesCategory,
			Description: "Create hydrated button",
			RunFunc:     HydratedRunFunc,
		})
}

func HydratedRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
	id, _ := args.Locals.Load("uid")
	button := util.CreateHydratedTemplateButton("Hello", "footer",
		util.GenerateHydratedUrlButton("url", "https://google.com"),
		util.GenerateHydratedCallButton("test", "+62 812 9798 0063"),
		util.GenerateHydratedQuickReplyButton("help", id, "!help"))
	return button
}
