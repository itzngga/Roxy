package util

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func SendReplyMessage(c *whatsmeow.Client, m *events.Message, text string) error {
	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        &text,
			ContextInfo: WithReply(m),
		},
	}
	_, err := c.SendMessage(m.Info.Chat, "", msg)
	return err
}
func ParseQuotedMessage(m *waProto.Message) *waProto.Message {
	if m.GetExtendedTextMessage().GetContextInfo() != nil {
		return m.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetImageMessage().GetContextInfo() != nil {
		return m.GetImageMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetVideoMessage().GetContextInfo() != nil {
		return m.GetVideoMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetDocumentMessage().GetContextInfo() != nil {
		return m.GetDocumentMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetAudioMessage().GetContextInfo() != nil {
		return m.GetAudioMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetStickerMessage().GetContextInfo() != nil {
		return m.GetStickerMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetButtonsMessage().GetContextInfo() != nil {
		return m.GetButtonsMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetGroupInviteMessage().GetContextInfo() != nil {
		return m.GetGroupInviteMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetProductMessage().GetContextInfo() != nil {
		return m.GetProductMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetListMessage().GetContextInfo() != nil {
		return m.GetListMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetTemplateMessage().GetContextInfo() != nil {
		return m.GetTemplateMessage().GetContextInfo().GetQuotedMessage()
	} else if m.GetContactMessage().GetContextInfo() != nil {
		return m.GetContactMessage().GetContextInfo().GetQuotedMessage()
	} else {
		return m
	}
}

func ParseQuotedBy(m *waProto.Message, str string) *waProto.Message {
	switch str {
	case "text":
		return m.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage()
	case "image":
		return m.GetImageMessage().GetContextInfo().GetQuotedMessage()
	case "video":
		return m.GetVideoMessage().GetContextInfo().GetQuotedMessage()
	case "sticker":
		return m.GetStickerMessage().GetContextInfo().GetQuotedMessage()
	case "document":
		return m.GetDocumentMessage().GetContextInfo().GetQuotedMessage()
	case "audio":
		return m.GetAudioMessage().GetContextInfo().GetQuotedMessage()
	default:
		return ParseQuotedMessage(m)
	}
}

func SendReplyText(m *events.Message, text string) *waProto.Message {
	return &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        &text,
			ContextInfo: WithReply(m),
		},
	}
}

func WithReply(m *events.Message) *waProto.ContextInfo {
	return &waProto.ContextInfo{
		StanzaId:      &m.Info.ID,
		Participant:   proto.String(m.Info.MessageSource.Sender.String()),
		QuotedMessage: m.Message,
	}
}

func SendInvalidCommand(m *events.Message) *waProto.Message {
	return &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String("Invalid command received"),
			ContextInfo: WithReply(m),
		},
	}
}
