package command

import (
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/itzngga/roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"strings"
)

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
	case []byte:
		mimetypeString := mimetype.Detect(value)
		if mimetypeString.Is("image/webp") {
			sticker, _ := runFunc.UploadStickerMessageFromBytes(value)
			message = &waProto.Message{
				StickerMessage: sticker,
			}
			sticker.ContextInfo = util.WithReply(runFunc.MessageEvent)
		} else if strings.Contains(mimetypeString.String(), "image") {
			image, _ := runFunc.UploadImageMessageFromBytes(value, "")
			message = &waProto.Message{
				ImageMessage: image,
			}
			image.ContextInfo = util.WithReply(runFunc.MessageEvent)
		} else if strings.Contains(mimetypeString.String(), "video") {
			video, _ := runFunc.UploadVideoMessageFromBytes(value, "")
			message = &waProto.Message{
				VideoMessage: video,
			}
			video.ContextInfo = util.WithReply(runFunc.MessageEvent)
		} else if strings.Contains(mimetypeString.String(), "audio") {
			audio, _ := runFunc.UploadAudioMessageFromBytes(value)
			message = &waProto.Message{
				AudioMessage: audio,
			}
			audio.ContextInfo = util.WithReply(runFunc.MessageEvent)
		} else {
			document, _ := runFunc.UploadDocumentMessageFromBytes(value, "", "document."+mimetypeString.Extension())
			message = &waProto.Message{
				DocumentMessage: document,
			}
			document.ContextInfo = util.WithReply(runFunc.MessageEvent)
		}
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
	case *waProto.ContactMessage:
		message = &waProto.Message{
			ContactMessage: value,
		}
		value.ContextInfo = util.WithReply(runFunc.MessageEvent)
	}

	_, err := runFunc.Client.SendMessage(context.Background(), runFunc.MessageEvent.Info.Chat, message)
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
				Text: &value,
			},
		}
	case []byte:
		mimetypeString := mimetype.Detect(value)
		if mimetypeString.Is("image/webp") {
			sticker, _ := runFunc.UploadStickerMessageFromBytes(value)
			message = &waProto.Message{
				StickerMessage: sticker,
			}
		} else if strings.Contains(mimetypeString.String(), "image") {
			image, _ := runFunc.UploadImageMessageFromBytes(value, "")
			message = &waProto.Message{
				ImageMessage: image,
			}
		} else if strings.Contains(mimetypeString.String(), "video") {
			video, _ := runFunc.UploadVideoMessageFromBytes(value, "")
			message = &waProto.Message{
				VideoMessage: video,
			}
		} else if strings.Contains(mimetypeString.String(), "audio") {
			audio, _ := runFunc.UploadAudioMessageFromBytes(value)
			message = &waProto.Message{
				AudioMessage: audio,
			}
		} else {
			document, _ := runFunc.UploadDocumentMessageFromBytes(value, "", "document."+mimetypeString.Extension())
			message = &waProto.Message{
				DocumentMessage: document,
			}
		}
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
	case *waProto.ContactMessage:
		message = &waProto.Message{
			ContactMessage: value,
		}
	}

	_, err := runFunc.Client.SendMessage(context.Background(), runFunc.MessageEvent.Info.Chat, message)
	if err != nil {
		fmt.Printf("error: sending message: %v\n", err)
		return
	}

	return
}
