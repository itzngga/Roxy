package command

import (
	"fmt"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func UrlImageButtonCommand() {
	AddCommand(&handler.Command{
		Name:        "url-button-image",
		Description: "Url Image button command",
		Category:    handler.UtilitiesCategory,
		RunFunc:     UrlImageButtonRunFunc,
	})
}

func UrlImageButtonRunFunc(c *whatsmeow.Client, m *events.Message, cmd *handler.Command) *waProto.Message {
	id := cmd.GetLocals("uid").(string)

	image, err := util.UploadImageFromUrl(c, "https://picsum.photos/id/870/1280/720?grayscale&blur=2", "Testing")
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
