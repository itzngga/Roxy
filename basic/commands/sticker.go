package commands

import (
	"fmt"
	"github.com/itzngga/goRoxy/basic/categories"
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/embed"
	"github.com/itzngga/goRoxy/util"
	"github.com/itzngga/goRoxy/util/cli"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	embed.Commands.Add(sticker)
}

var sticker = &command.Command{
	Name:        "sticker",
	Aliases:     []string{"s", "stiker"},
	Category:    categories.CommonCategory,
	Description: "Create sticker from image or video",
	RunFunc: func(c *whatsmeow.Client, params *command.RunFuncParams) *waProto.Message {
		if params.Message.GetImageMessage() != nil {
			return StickerImage(c, params.Event, params.Message.GetImageMessage())
		} else if util.ParseQuotedMessage(params.Message).GetImageMessage() != nil {
			return StickerImage(c, params.Event, util.ParseQuotedMessage(params.Message).GetImageMessage())
		} else if params.Message.GetVideoMessage() != nil {
			return StickerVideo(c, params.Event, params.Message.GetVideoMessage())
		} else if util.ParseQuotedMessage(params.Message).GetVideoMessage() != nil {
			return StickerVideo(c, params.Event, util.ParseQuotedMessage(params.Message).GetVideoMessage())
		}
		return util.GenerateReplyMessage(params.Event, "Invalid")
	},
}

func StickerVideo(c *whatsmeow.Client, event *events.Message, video *waProto.VideoMessage) *waProto.Message {
	data, err := c.Download(video)
	if err != nil {
		fmt.Printf("Failed to download video: %v\n", err)
	}

	var qValue int
	switch dataLen := len(data); {
	case dataLen < 300000:
		qValue = 20
	case dataLen < 400000:
		qValue = 10
	default:
		qValue = 5
	}

	resultData := cli.FfmpegPipeline(data,
		"-y", "-hide_banner", "-loglevel", "panic",
		"-i", "pipe:0",
		"-filter:v", "fps=fps=15",
		"-compression_level", "0",
		"-q:v", fmt.Sprintf("%d", qValue),
		"-loop", "0",
		"-preset", "picture",
		"-an", "-vsync", "0",
		"-s", "512:512",
		"-f", "webp",
		"pipe:1",
	)

	util.SendReplyMessage(c, event, resultData)
	return nil
}
func StickerImage(c *whatsmeow.Client, event *events.Message, img *waProto.ImageMessage) *waProto.Message {
	data, err := c.Download(img)
	if err != nil {
		fmt.Printf("Failed to download image: %v\n", err)
	}

	resultData := cli.FfmpegPipeline(data,
		"-y", "-hide_banner", "-loglevel", "panic",
		"-i", "pipe:0",
		"-f", "webp",
		"-s", "512:512",
		"-preset", "picture",
		"pipe:1",
	)

	util.SendReplyMessage(c, event, resultData)
	return nil
}
