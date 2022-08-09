package command

import (
	"fmt"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func ImageButtonCommand() {
	AddCommand(&handler.Command{
		Name:        "button-image",
		Description: "Image button command",
		Category:    handler.UtilitiesCategory,
		RunFunc:     ImageButtonRunFunc,
	})
}

func ImageButtonRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
	id := args.Cmd.GetLocals("uid").(string)

	image, err := util.UploadImageMessageFromPath(c, "temp/example.png", "Testing")
	if err != nil {
		fmt.Printf("\nError uploading image :  %v", err)
		return nil
		// remember to return nil when error is returned
	}

	return util.CreateImageButton("test", "footer",
		&waProto.ButtonsMessage_ImageMessage{
			ImageMessage: image,
		},
		util.GenerateButton(id, "!help", "Help"),
	)
}
