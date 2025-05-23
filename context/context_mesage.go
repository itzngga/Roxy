package context

import (
	"fmt"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/itzngga/Roxy/util"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/proto/waHistorySync"
	waTypes "go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// RevokeMessage revoke message from given jid and client message id
func (context *Ctx) RevokeMessage(jid waTypes.JID, messageId waTypes.MessageID) {
	_, _ = context.Methods().SendMessage(jid, context.Client().BuildRevoke(jid, waTypes.EmptyJID, messageId))
}

// ByteToMessage convert byte to whatsmeow message object
func (context *Ctx) ByteToMessage(value []byte, withReply bool, caption string) *waE2E.Message {
	var message *waE2E.Message
	mimetypeString := mimetype.Detect(value)
	if mimetypeString.Is("image/webp") {
		sticker, _ := context.UploadStickerMessageFromBytes(value)
		message = &waE2E.Message{
			StickerMessage: sticker,
		}
		if withReply {
			sticker.ContextInfo = util.WithReply(context.MessageEvent())
		}
	} else if strings.Contains(mimetypeString.String(), "image") {
		image, _ := context.UploadImageMessageFromBytes(value, caption)
		message = &waE2E.Message{
			ImageMessage: image,
		}
		if withReply {
			image.ContextInfo = util.WithReply(context.MessageEvent())
		}
	} else if strings.Contains(mimetypeString.String(), "video") {
		video, _ := context.UploadVideoMessageFromBytes(value, caption)
		message = &waE2E.Message{
			VideoMessage: video,
		}
		if withReply {
			video.ContextInfo = util.WithReply(context.MessageEvent())
		}
	} else if strings.Contains(mimetypeString.String(), "audio") {
		audio, _ := context.UploadAudioMessageFromBytes(value)
		message = &waE2E.Message{
			AudioMessage: audio,
		}
		if withReply {
			audio.ContextInfo = util.WithReply(context.MessageEvent())
		}
	} else {
		document, _ := context.UploadDocumentMessageFromBytes(value, caption, "document."+mimetypeString.Extension())
		message = &waE2E.Message{
			DocumentMessage: document,
		}
		if withReply {
			document.ContextInfo = util.WithReply(context.MessageEvent())
		}
	}
	return message
}

// SendReplyMessage send reply message in current chat
func (context *Ctx) SendReplyMessage(obj any) {
	var message *waE2E.Message
	switch value := obj.(type) {
	case string:
		message = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        &value,
				ContextInfo: util.WithReply(context.MessageEvent()),
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
				message = context.ByteToMessage(byteResult, true, caption)
			}
		}
	case []byte:
		message = context.ByteToMessage(value, true, "")
	case *waHistorySync.Conversation:
		a := value.String()
		message = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        &a,
				ContextInfo: util.WithReply(context.MessageEvent()),
			},
		}
	case *waE2E.ImageMessage:
		message = &waE2E.Message{
			ImageMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.ExtendedTextMessage:
		message = &waE2E.Message{
			ExtendedTextMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.DocumentMessage:
		message = &waE2E.Message{
			DocumentMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.VideoMessage:
		message = &waE2E.Message{
			VideoMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.AudioMessage:
		message = &waE2E.Message{
			AudioMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.StickerMessage:
		message = &waE2E.Message{
			StickerMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	// case *waProto.ButtonsMessage:
	//	message = &waProto.Message{
	//		ButtonsMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.GroupInviteMessage:
		message = &waE2E.Message{
			GroupInviteMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.ProductMessage:
		message = &waE2E.Message{
			ProductMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	// case *waProto.ListMessage:
	//	message = &waProto.Message{
	//		ListMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	// case *waProto.TemplateMessage:
	//	message = &waProto.Message{
	//		TemplateMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.Message:
		message = value
	case *waE2E.ContactMessage:
		message = &waE2E.Message{
			ContactMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	}

	context.Methods().SendMessage(context.MessageEvent().Info.Chat, message)
}

// GenerateReplyMessage generate reply message to whatsmeow message object
func (context *Ctx) GenerateReplyMessage(obj any) *waE2E.Message {
	var message *waE2E.Message
	switch value := obj.(type) {
	case string:
		message = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        &value,
				ContextInfo: util.WithReply(context.MessageEvent()),
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
				message = context.ByteToMessage(byteResult, true, caption)
			}
		}
	case []byte:
		message = context.ByteToMessage(value, true, "")
	case *waHistorySync.Conversation:
		a := value.String()
		message = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        &a,
				ContextInfo: util.WithReply(context.MessageEvent()),
			},
		}
	case *waE2E.ImageMessage:
		message = &waE2E.Message{
			ImageMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.ExtendedTextMessage:
		message = &waE2E.Message{
			ExtendedTextMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.DocumentMessage:
		message = &waE2E.Message{
			DocumentMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.VideoMessage:
		message = &waE2E.Message{
			VideoMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.AudioMessage:
		message = &waE2E.Message{
			AudioMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.StickerMessage:
		message = &waE2E.Message{
			StickerMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	// case *waProto.ButtonsMessage:
	//	message = &waProto.Message{
	//		ButtonsMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.GroupInviteMessage:
		message = &waE2E.Message{
			GroupInviteMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.ProductMessage:
		message = &waE2E.Message{
			ProductMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	// case *waProto.ListMessage:
	//	message = &waProto.Message{
	//		ListMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	// case *waProto.TemplateMessage:
	//	message = &waProto.Message{
	//		TemplateMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waE2E.Message:
		message = value
	case *waE2E.ContactMessage:
		message = &waE2E.Message{
			ContactMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	}

	return message
}

func (context *Ctx) SendMessage(obj any) {
	var message *waE2E.Message
	switch value := obj.(type) {
	case string:
		message = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        &value,
				ContextInfo: util.WithReply(context.MessageEvent()),
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
				message = context.ByteToMessage(byteResult, false, caption)
			}
		}
	case []byte:
		message = context.ByteToMessage(value, false, "")
	case *waHistorySync.Conversation:
		a := value.String()
		message = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text: &a,
			},
		}
	case *waE2E.ImageMessage:
		message = &waE2E.Message{
			ImageMessage: value,
		}
	case *waE2E.ExtendedTextMessage:
		message = &waE2E.Message{
			ExtendedTextMessage: value,
		}
	case *waE2E.DocumentMessage:
		message = &waE2E.Message{
			DocumentMessage: value,
		}
	case *waE2E.VideoMessage:
		message = &waE2E.Message{
			VideoMessage: value,
		}
	case *waE2E.AudioMessage:
		message = &waE2E.Message{
			AudioMessage: value,
		}
	case *waE2E.StickerMessage:
		message = &waE2E.Message{
			StickerMessage: value,
		}
	// case *waProto.ButtonsMessage:
	//	message = &waProto.Message{
	//		ButtonsMessage: value,
	//	}
	case *waE2E.GroupInviteMessage:
		message = &waE2E.Message{
			GroupInviteMessage: value,
		}
	case *waE2E.ProductMessage:
		message = &waE2E.Message{
			ProductMessage: value,
		}
	// case *waProto.ListMessage:
	//	message = &waProto.Message{
	//		ListMessage: value,
	//	}
	// case *waProto.TemplateMessage:
	//	message = &waProto.Message{
	//		TemplateMessage: value,
	//	}
	case *waE2E.Message:
		message = value
	case *waE2E.ContactMessage:
		message = &waE2E.Message{
			ContactMessage: value,
		}
	}

	_, _ = context.Methods().SendMessage(context.MessageEvent().Info.Chat, message)
}

// EditMessageText edit current text message to given text
func (context *Ctx) EditMessageText(to string) error {
	msgKey := &waCommon.MessageKey{
		RemoteJID: proto.String(context.ChatJID().String()),
		FromMe:    proto.Bool(true),
		ID:        proto.String(context.MessageInfo().ID),
	}
	if context.Message().GetExtendedTextMessage() != nil {
		context.Message().ExtendedTextMessage.Text = proto.String(to)
	} else if context.Message().GetConversation() != "" {
		context.Message().Conversation = proto.String(to)
	} else {
		return fmt.Errorf("error: invalid message type")
	}

	message := &waE2E.Message{
		ProtocolMessage: &waE2E.ProtocolMessage{
			Key:           msgKey,
			Type:          (*waE2E.ProtocolMessage_Type)(proto.Int32(14)),
			EditedMessage: context.Message(),
		},
	}

	_, err := context.Methods().SendMessage(context.MessageEvent().Info.Chat, message)
	return err
}

// SendEmoji send emoji to current text message
func (context *Ctx) SendEmoji(emoji string) {
	context.Methods().SendEmojiMessage(context.MessageEvent(), emoji)
}
