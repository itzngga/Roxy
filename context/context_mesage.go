package context

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"strings"
)

// RevokeMessage revoke message from given jid and client message id
func (context *Ctx) RevokeMessage(jid waTypes.JID, messageId waTypes.MessageID) {
	_, _ = context.Methods().SendMessage(jid, context.Client().BuildRevoke(jid, waTypes.EmptyJID, messageId))
	return
}

// ByteToMessage convert byte to whatsmeow message object
func (context *Ctx) ByteToMessage(value []byte, withReply bool, caption string) *waProto.Message {
	var message *waProto.Message
	mimetypeString := mimetype.Detect(value)
	if mimetypeString.Is("image/webp") {
		sticker, _ := context.UploadStickerMessageFromBytes(value)
		message = &waProto.Message{
			StickerMessage: sticker,
		}
		if withReply {
			sticker.ContextInfo = util.WithReply(context.MessageEvent())
		}
	} else if strings.Contains(mimetypeString.String(), "image") {
		image, _ := context.UploadImageMessageFromBytes(value, caption)
		message = &waProto.Message{
			ImageMessage: image,
		}
		if withReply {
			image.ContextInfo = util.WithReply(context.MessageEvent())
		}
	} else if strings.Contains(mimetypeString.String(), "video") {
		video, _ := context.UploadVideoMessageFromBytes(value, caption)
		message = &waProto.Message{
			VideoMessage: video,
		}
		if withReply {
			video.ContextInfo = util.WithReply(context.MessageEvent())
		}
	} else if strings.Contains(mimetypeString.String(), "audio") {
		audio, _ := context.UploadAudioMessageFromBytes(value)
		message = &waProto.Message{
			AudioMessage: audio,
		}
		if withReply {
			audio.ContextInfo = util.WithReply(context.MessageEvent())
		}
	} else {
		document, _ := context.UploadDocumentMessageFromBytes(value, caption, "document."+mimetypeString.Extension())
		message = &waProto.Message{
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
	var message *waProto.Message
	switch value := obj.(type) {
	case string:
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
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
	case *waProto.Conversation:
		a := value.String()
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &a,
				ContextInfo: util.WithReply(context.MessageEvent()),
			},
		}
	case *waProto.ImageMessage:
		message = &waProto.Message{
			ImageMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.ExtendedTextMessage:
		message = &waProto.Message{
			ExtendedTextMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.DocumentMessage:
		message = &waProto.Message{
			DocumentMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.VideoMessage:
		message = &waProto.Message{
			VideoMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.AudioMessage:
		message = &waProto.Message{
			AudioMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.StickerMessage:
		message = &waProto.Message{
			StickerMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	//case *waProto.ButtonsMessage:
	//	message = &waProto.Message{
	//		ButtonsMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.GroupInviteMessage:
		message = &waProto.Message{
			GroupInviteMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.ProductMessage:
		message = &waProto.Message{
			ProductMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	//case *waProto.ListMessage:
	//	message = &waProto.Message{
	//		ListMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	//case *waProto.TemplateMessage:
	//	message = &waProto.Message{
	//		TemplateMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.Message:
		message = value
	case *waProto.ContactMessage:
		message = &waProto.Message{
			ContactMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	}

	_, _ = context.Methods().SendMessage(context.MessageEvent().Info.Chat, message)
	return
}

// GenerateReplyMessage generate reply message to whatsmeow message object
func (context *Ctx) GenerateReplyMessage(obj any) *waProto.Message {
	var message *waProto.Message
	switch value := obj.(type) {
	case string:
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
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
	case *waProto.Conversation:
		a := value.String()
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &a,
				ContextInfo: util.WithReply(context.MessageEvent()),
			},
		}
	case *waProto.ImageMessage:
		message = &waProto.Message{
			ImageMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.ExtendedTextMessage:
		message = &waProto.Message{
			ExtendedTextMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.DocumentMessage:
		message = &waProto.Message{
			DocumentMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.VideoMessage:
		message = &waProto.Message{
			VideoMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.AudioMessage:
		message = &waProto.Message{
			AudioMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.StickerMessage:
		message = &waProto.Message{
			StickerMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	//case *waProto.ButtonsMessage:
	//	message = &waProto.Message{
	//		ButtonsMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.GroupInviteMessage:
		message = &waProto.Message{
			GroupInviteMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.ProductMessage:
		message = &waProto.Message{
			ProductMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	//case *waProto.ListMessage:
	//	message = &waProto.Message{
	//		ListMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	//case *waProto.TemplateMessage:
	//	message = &waProto.Message{
	//		TemplateMessage: value,
	//	}
	//	value.ContextInfo = util.WithReply(context.MessageEvent())
	case *waProto.Message:
		message = value
	case *waProto.ContactMessage:
		message = &waProto.Message{
			ContactMessage: value,
		}
		value.ContextInfo = util.WithReply(context.MessageEvent())
	}

	return message
}

func (context *Ctx) SendMessage(obj any) {
	var message *waProto.Message
	switch value := obj.(type) {
	case string:
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
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

	_, _ = context.Methods().SendMessage(context.MessageEvent().Info.Chat, message)
	return
}

// EditMessageText edit current text message to given text
func (context *Ctx) EditMessageText(to string) error {
	msgKey := &waProto.MessageKey{
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

	message := &waProto.Message{
		ProtocolMessage: &waProto.ProtocolMessage{
			Key:           msgKey,
			Type:          (*waProto.ProtocolMessage_Type)(proto.Int32(14)),
			EditedMessage: context.Message(),
		},
	}

	_, err := context.Methods().SendMessage(context.MessageEvent().Info.Chat, message)
	return err
}

// SendEmoji send emoji to current text message
func (context *Ctx) SendEmoji(emoji string) {
	context.Methods().SendEmojiMessage(context.MessageEvent(), emoji)
	return
}
