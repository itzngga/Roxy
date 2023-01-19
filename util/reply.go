package util

import (
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"strings"
)

func SendReplyMessage(c *whatsmeow.Client, event *events.Message, obj any) {
	var message *waProto.Message
	switch value := obj.(type) {
	case string:
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &value,
				ContextInfo: WithReply(event),
			},
		}
	case []byte:
		mimetypeString := mimetype.Detect(value)
		if mimetypeString.Is("image/webp") {
			sticker, _ := UploadStickerMessageFromBytes(c, event, value)
			message = &waProto.Message{
				StickerMessage: sticker,
			}
			sticker.ContextInfo = WithReply(event)
		} else if strings.Contains(mimetypeString.String(), "image") {
			image, _ := UploadImageMessageFromBytes(c, event, value, "")
			message = &waProto.Message{
				ImageMessage: image,
			}
			image.ContextInfo = WithReply(event)
		} else if strings.Contains(mimetypeString.String(), "video") {
			video, _ := UploadVideoMessageFromBytes(c, event, value, "")
			message = &waProto.Message{
				VideoMessage: video,
			}
			video.ContextInfo = WithReply(event)
		} else if strings.Contains(mimetypeString.String(), "audio") {
			audio, _ := UploadAudioMessageFromBytes(c, value)
			message = &waProto.Message{
				AudioMessage: audio,
			}
			audio.ContextInfo = WithReply(event)
		} else {
			document, _ := UploadDocumentMessageFromBytes(c, value, "", "document."+mimetypeString.Extension())
			message = &waProto.Message{
				DocumentMessage: document,
			}
			document.ContextInfo = WithReply(event)
		}
	case *waProto.Conversation:
		a := value.String()
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &a,
				ContextInfo: WithReply(event),
			},
		}
	case *waProto.ImageMessage:
		message = &waProto.Message{
			ImageMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.ExtendedTextMessage:
		message = &waProto.Message{
			ExtendedTextMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.DocumentMessage:
		message = &waProto.Message{
			DocumentMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.VideoMessage:
		message = &waProto.Message{
			VideoMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.AudioMessage:
		message = &waProto.Message{
			AudioMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.StickerMessage:
		message = &waProto.Message{
			StickerMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.ButtonsMessage:
		message = &waProto.Message{
			ButtonsMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.GroupInviteMessage:
		message = &waProto.Message{
			GroupInviteMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.ProductMessage:
		message = &waProto.Message{
			ProductMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.ListMessage:
		message = &waProto.Message{
			ListMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.TemplateMessage:
		message = &waProto.Message{
			TemplateMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.ContactMessage:
		message = &waProto.Message{
			ContactMessage: value,
		}
		value.ContextInfo = WithReply(event)
	}

	_, err := c.SendMessage(context.Background(), event.Info.Chat, message)
	if err != nil {
		fmt.Printf("error: sending message: %v\n", err)
	}
}

func GenerateReplyMessage(event *events.Message, obj any) *waProto.Message {
	var message *waProto.Message
	switch value := obj.(type) {
	case string:
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &value,
				ContextInfo: WithReply(event),
			},
		}
	case *waProto.Conversation:
		a := value.String()
		message = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &a,
				ContextInfo: WithReply(event),
			},
		}
	case *waProto.ImageMessage:
		message = &waProto.Message{
			ImageMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.ExtendedTextMessage:
		message = &waProto.Message{
			ExtendedTextMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.DocumentMessage:
		message = &waProto.Message{
			DocumentMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.VideoMessage:
		message = &waProto.Message{
			VideoMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.AudioMessage:
		message = &waProto.Message{
			AudioMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.StickerMessage:
		message = &waProto.Message{
			StickerMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.ButtonsMessage:
		message = &waProto.Message{
			ButtonsMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.GroupInviteMessage:
		message = &waProto.Message{
			GroupInviteMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.ProductMessage:
		message = &waProto.Message{
			ProductMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.ListMessage:
		message = &waProto.Message{
			ListMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.TemplateMessage:
		message = &waProto.Message{
			TemplateMessage: value,
		}
		value.ContextInfo = WithReply(event)
	case *waProto.ContactMessage:
		message = &waProto.Message{
			ContactMessage: value,
		}
		value.ContextInfo = WithReply(event)
	}

	return message
}
