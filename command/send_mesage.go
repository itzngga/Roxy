package command

import (
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/itzngga/Roxy/types"
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"strings"
	"time"
)

func (runFunc *RunFuncContext) ByteToMessage(value []byte, withReply bool, caption string) *waProto.Message {
	var message *waProto.Message
	mimetypeString := mimetype.Detect(value)
	if mimetypeString.Is("image/webp") {
		sticker, _ := runFunc.UploadStickerMessageFromBytes(value)
		message = &waProto.Message{
			StickerMessage: sticker,
		}
		if withReply {
			sticker.ContextInfo = util.WithReply(runFunc.MessageEvent)
		}
	} else if strings.Contains(mimetypeString.String(), "image") {
		image, _ := runFunc.UploadImageMessageFromBytes(value, caption)
		message = &waProto.Message{
			ImageMessage: image,
		}
		if withReply {
			image.ContextInfo = util.WithReply(runFunc.MessageEvent)
		}
	} else if strings.Contains(mimetypeString.String(), "video") {
		video, _ := runFunc.UploadVideoMessageFromBytes(value, caption)
		message = &waProto.Message{
			VideoMessage: video,
		}
		if withReply {
			video.ContextInfo = util.WithReply(runFunc.MessageEvent)
		}
	} else if strings.Contains(mimetypeString.String(), "audio") {
		audio, _ := runFunc.UploadAudioMessageFromBytes(value)
		message = &waProto.Message{
			AudioMessage: audio,
		}
		if withReply {
			audio.ContextInfo = util.WithReply(runFunc.MessageEvent)
		}
	} else {
		document, _ := runFunc.UploadDocumentMessageFromBytes(value, caption, "document."+mimetypeString.Extension())
		message = &waProto.Message{
			DocumentMessage: document,
		}
		if withReply {
			document.ContextInfo = util.WithReply(runFunc.MessageEvent)
		}
	}
	return message
}

func (runFunc *RunFuncContext) SendReplyMessage(obj any) {
	var message *waProto.Message
	switch value := obj.(type) {
	case string:
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &value,
				ContextInfo: util.WithReply(runFunc.MessageEvent),
			},
		}
	case map[string]string:
		for url, caption := range value {
			if util.IsValidUrl(url) {
				byteResult, err := util.DoHTTPRequest("GET", url)
				if err != nil {
					fmt.Println("error: " + err.Error())
					return
				}
				message = runFunc.ByteToMessage(byteResult, true, caption)
			}
		}
	case []byte:
		message = runFunc.ByteToMessage(value, true, "")
	case *waProto.Conversation:
		a := value.String()
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &a,
				ContextInfo: util.WithReply(runFunc.MessageEvent),
			},
		}
	case *waProto.ImageMessage:
		message = &waProto.Message{
			ImageMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.ExtendedTextMessage:
		message = &waProto.Message{
			ExtendedTextMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.DocumentMessage:
		message = &waProto.Message{
			DocumentMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.VideoMessage:
		message = &waProto.Message{
			VideoMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.AudioMessage:
		message = &waProto.Message{
			AudioMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.StickerMessage:
		message = &waProto.Message{
			StickerMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	//case *waProto.ButtonsMessage:
	//	message = &waProto.Message{
	//		ButtonsMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.GroupInviteMessage:
		message = &waProto.Message{
			GroupInviteMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.ProductMessage:
		message = &waProto.Message{
			ProductMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	//case *waProto.ListMessage:
	//	message = &waProto.Message{
	//		ListMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	//case *waProto.TemplateMessage:
	//	message = &waProto.Message{
	//		TemplateMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.Message:
		message = value
	case *waProto.ContactMessage:
		message = &waProto.Message{
			ContactMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	}

	ctx, cancel := context.WithTimeout(context.Background(), runFunc.Options.SendMessageTimeout)
	defer cancel()

	_, err := runFunc.Client.SendMessage(ctx, runFunc.MessageEvent.Info.Chat, message)
	if err != nil {
		fmt.Printf("error: sending message: %v\n", err)
	}
}

func (runFunc *RunFuncContext) GenerateReplyMessage(obj any) *waProto.Message {
	var message *waProto.Message
	switch value := obj.(type) {
	case string:
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &value,
				ContextInfo: util.WithReply(runFunc.MessageEvent),
			},
		}
	case map[string]string:
		for url, caption := range value {
			if util.IsValidUrl(url) {
				byteResult, err := util.DoHTTPRequest("GET", url)
				if err != nil {
					fmt.Println("error: " + err.Error())
					return nil
				}
				message = runFunc.ByteToMessage(byteResult, true, caption)
			}
		}
	case []byte:
		message = runFunc.ByteToMessage(value, true, "")
	case *waProto.Conversation:
		a := value.String()
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &a,
				ContextInfo: util.WithReply(runFunc.MessageEvent),
			},
		}
	case *waProto.ImageMessage:
		message = &waProto.Message{
			ImageMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.ExtendedTextMessage:
		message = &waProto.Message{
			ExtendedTextMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.DocumentMessage:
		message = &waProto.Message{
			DocumentMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.VideoMessage:
		message = &waProto.Message{
			VideoMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.AudioMessage:
		message = &waProto.Message{
			AudioMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.StickerMessage:
		message = &waProto.Message{
			StickerMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	//case *waProto.ButtonsMessage:
	//	message = &waProto.Message{
	//		ButtonsMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.GroupInviteMessage:
		message = &waProto.Message{
			GroupInviteMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.ProductMessage:
		message = &waProto.Message{
			ProductMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	//case *waProto.ListMessage:
	//	message = &waProto.Message{
	//		ListMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	//case *waProto.TemplateMessage:
	//	message = &waProto.Message{
	//		TemplateMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	case *waProto.Message:
		message = value
	case *waProto.ContactMessage:
		message = &waProto.Message{
			ContactMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	}

	return message
}

func (runFunc *RunFuncContext) SendMessage(obj any) {
	var message *waProto.Message
	switch value := obj.(type) {
	case string:
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &value,
				ContextInfo: util.WithReply(runFunc.MessageEvent),
			},
		}
	case map[string]string:
		for url, caption := range value {
			if util.IsValidUrl(url) {
				byteResult, err := util.DoHTTPRequest("GET", url)
				if err != nil {
					fmt.Println("error: " + err.Error())
					return
				}
				message = runFunc.ByteToMessage(byteResult, false, caption)
			}
		}
	case []byte:
		message = runFunc.ByteToMessage(value, false, "")
	case *waProto.Conversation:
		a := value.String()
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text: &a,
			},
		}
	case *waProto.ImageMessage:
		message = &waProto.Message{
			ImageMessage: value,
		}
	case *waProto.ExtendedTextMessage:
		message = &waProto.Message{
			ExtendedTextMessage: value,
		}
	case *waProto.DocumentMessage:
		message = &waProto.Message{
			DocumentMessage: value,
		}
	case *waProto.VideoMessage:
		message = &waProto.Message{
			VideoMessage: value,
		}
	case *waProto.AudioMessage:
		message = &waProto.Message{
			AudioMessage: value,
		}
	case *waProto.StickerMessage:
		message = &waProto.Message{
			StickerMessage: value,
		}
	//case *waProto.ButtonsMessage:
	//	message = &waProto.Message{
	//		ButtonsMessage: value,
	//	}
	case *waProto.GroupInviteMessage:
		message = &waProto.Message{
			GroupInviteMessage: value,
		}
	case *waProto.ProductMessage:
		message = &waProto.Message{
			ProductMessage: value,
		}
	//case *waProto.ListMessage:
	//	message = &waProto.Message{
	//		ListMessage: value,
	//	}
	//case *waProto.TemplateMessage:
	//	message = &waProto.Message{
	//		TemplateMessage: value,
	//	}
	case *waProto.Message:
		message = value
	case *waProto.ContactMessage:
		message = &waProto.Message{
			ContactMessage: value,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), runFunc.Options.SendMessageTimeout)
	defer cancel()

	_, err := runFunc.Client.SendMessage(ctx, runFunc.MessageEvent.Info.Chat, message)
	if err != nil {
		fmt.Printf("error: sending message: %v\n", err)
		return
	}

	return
}

func (runFunc *RunFuncContext) EditMessageText(to string) error {
	msgKey := &waProto.MessageKey{
		RemoteJid: types.String(runFunc.MessageInfo.Chat.String()),
		FromMe:    types.Bool(true),
		Id:        types.String(runFunc.MessageInfo.ID),
	}
	if runFunc.Message.GetExtendedTextMessage() != nil {
		runFunc.Message.ExtendedTextMessage.Text = types.String(to)
	} else if runFunc.Message.GetConversation() != "" {
		runFunc.Message.Conversation = types.String(to)
	} else {
		return fmt.Errorf("error: invalid message type")
	}

	message := &waProto.Message{
		ProtocolMessage: &waProto.ProtocolMessage{
			Key:           msgKey,
			Type:          (*waProto.ProtocolMessage_Type)(types.Int32(14)),
			EditedMessage: runFunc.Message,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*runFunc.Options.SendMessageTimeout)
	defer cancel()

	_, err := runFunc.Client.SendMessage(ctx, runFunc.MessageEvent.Info.Chat, message)
	if err != nil {
		return fmt.Errorf("error: sending message: %v\n", err)
	}

	return nil
}
