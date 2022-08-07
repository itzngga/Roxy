package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
	"mime"

	"net/http"
	"os"
	"os/exec"
)

func StickerCommand() {
	AddCommand(
		&handler.Command{
			Name:        "sticker",
			Aliases:     []string{"stkr", "stiker"},
			Category:    handler.UtilitiesCategory,
			Description: "Create sticker from image or video",
			RunFunc:     StickerRunFunc,
		})
}

func StickerRunFunc(c *whatsmeow.Client, m *events.Message, cmd *handler.Command) *waProto.Message {
	if m.Message.GetImageMessage() != nil {
		return StickerImage(c, m, m.Message.GetImageMessage())
	} else if util.ParseQuotedMessage(m.Message).GetImageMessage() != nil {
		return StickerImage(c, m, util.ParseQuotedMessage(m.Message).GetImageMessage())
	} else if m.Message.GetVideoMessage() != nil {
		return StickerVideo(c, m, m.Message.GetVideoMessage())
	} else if util.ParseQuotedMessage(m.Message).GetVideoMessage() != nil {
		return StickerVideo(c, m, util.ParseQuotedMessage(m.Message).GetVideoMessage())
	}
	return util.SendReplyText(m, "Invalid")
}
func StickerVideo(c *whatsmeow.Client, m *events.Message, video *waProto.VideoMessage) *waProto.Message {
	data, err := c.Download(video)
	if err != nil {
		fmt.Printf("Failed to download video: %v\n", err)
	}
	exts, _ := mime.ExtensionsByType(video.GetMimetype())
	RawPath := fmt.Sprintf("temp/%s%s", m.Info.ID, exts[0])
	ConvertedPath := fmt.Sprintf("temp/%s%s", m.Info.ID, ".webp")
	err = os.WriteFile(RawPath, data, 0600)
	if err != nil {
		fmt.Printf("Failed to save video: %v", err)
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

	commandString := fmt.Sprintf("ffmpeg -i %s -vcodec libwebp -filter:v fps=fps=15 -compression_level 0 -q:v %d -loop 0 -preset picture -an -vsync 0 -s 512:512 %s", RawPath, qValue, ConvertedPath)
	cmd := exec.Command("bash", "-c", commandString)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		fmt.Println(outb.String(), "*****", errb.String())
		fmt.Printf("Failed to Convert Video to WebP %s", err)
	}

	data, err = os.ReadFile(ConvertedPath)
	if err != nil {
		fmt.Printf("Failed to read %s: %s\n", ConvertedPath, err)
	}

	uploaded, err := c.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
	}
	defer os.Remove(RawPath)
	defer os.Remove(ConvertedPath)

	return &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			IsAnimated:    proto.Bool(true),
			ContextInfo:   util.WithReply(m),
		},
	}
}
func StickerImage(c *whatsmeow.Client, m *events.Message, img *waProto.ImageMessage) *waProto.Message {
	//vips.Startup(nil)
	//defer vips.Shutdown()

	data, err := c.Download(img)
	if err != nil {
		fmt.Printf("Failed to download image: %v\n", err)
	}
	exts, _ := mime.ExtensionsByType(img.GetMimetype())
	RawPath := fmt.Sprintf("temp/%s%s", m.Info.ID, exts[0])
	ConvertedPath := fmt.Sprintf("temp/%s%s", m.Info.ID, ".webp")
	err = os.WriteFile(RawPath, data, 0600)
	if err != nil {
		fmt.Printf("Failed to save image: %v", err)
	}
	cmd := exec.Command("cwebp", RawPath, "-resize", "0", "600", "-o", ConvertedPath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Failed to Convert Image to WebP")
	}

	// libvips code
	//input, err := vips.NewImageFromBuffer(data)
	//defer input.Close()
	//err = input.OptimizeICCProfile()
	//if err != nil {
	//	fmt.Println("Failed to Convert Image to WebP")
	//}
	//out := vips.NewWebpExportParams()
	//out.Lossless = true
	//data, _, err = input.ExportWebp(out)
	//
	//if err != nil {
	//	fmt.Println("Failed to Convert Image to WebP")
	//}
	data, err = os.ReadFile(ConvertedPath)
	if err != nil {
		fmt.Printf("Failed to read %s: %s\n", ConvertedPath, err)
	}

	//Upload WebP
	uploaded, err := c.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
	}
	defer os.Remove(RawPath)
	defer os.Remove(ConvertedPath)

	// Send WebP as sticker
	return &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo:   util.WithReply(m),
		},
	}

}
