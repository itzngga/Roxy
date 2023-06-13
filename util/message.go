package util

import (
	"context"
	"fmt"
	"github.com/itzngga/roxy/types"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"strings"
)

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

func ParseQuotedRemoteJid(evm *events.Message) *string {
	m := evm.Message
	if m.GetExtendedTextMessage().GetContextInfo() != nil {
		return m.GetExtendedTextMessage().GetContextInfo().Participant
	} else if m.GetImageMessage().GetContextInfo() != nil {
		return m.GetImageMessage().GetContextInfo().Participant
	} else if m.GetVideoMessage().GetContextInfo() != nil {
		return m.GetVideoMessage().GetContextInfo().Participant
	} else if m.GetDocumentMessage().GetContextInfo() != nil {
		return m.GetDocumentMessage().GetContextInfo().Participant
	} else if m.GetAudioMessage().GetContextInfo() != nil {
		return m.GetAudioMessage().GetContextInfo().Participant
	} else if m.GetStickerMessage().GetContextInfo() != nil {
		return m.GetStickerMessage().GetContextInfo().Participant
	} else if m.GetButtonsMessage().GetContextInfo() != nil {
		return m.GetButtonsMessage().GetContextInfo().Participant
	} else if m.GetGroupInviteMessage().GetContextInfo() != nil {
		return m.GetGroupInviteMessage().GetContextInfo().Participant
	} else if m.GetProductMessage().GetContextInfo() != nil {
		return m.GetProductMessage().GetContextInfo().Participant
	} else if m.GetListMessage().GetContextInfo() != nil {
		return m.GetListMessage().GetContextInfo().Participant
	} else if m.GetTemplateMessage().GetContextInfo() != nil {
		return m.GetTemplateMessage().GetContextInfo().Participant
	} else if m.GetContactMessage().GetContextInfo() != nil {
		return m.GetContactMessage().GetContextInfo().Participant
	} else {
		return &evm.Info.Chat.User
	}
}

func ParseQuotedMessageId(evm *events.Message) *string {
	m := evm.Message
	if m.GetExtendedTextMessage().GetContextInfo() != nil {
		return m.GetExtendedTextMessage().GetContextInfo().StanzaId
	} else if m.GetImageMessage().GetContextInfo() != nil {
		return m.GetImageMessage().GetContextInfo().StanzaId
	} else if m.GetVideoMessage().GetContextInfo() != nil {
		return m.GetVideoMessage().GetContextInfo().StanzaId
	} else if m.GetDocumentMessage().GetContextInfo() != nil {
		return m.GetDocumentMessage().GetContextInfo().StanzaId
	} else if m.GetAudioMessage().GetContextInfo() != nil {
		return m.GetAudioMessage().GetContextInfo().StanzaId
	} else if m.GetStickerMessage().GetContextInfo() != nil {
		return m.GetStickerMessage().GetContextInfo().StanzaId
	} else if m.GetButtonsMessage().GetContextInfo() != nil {
		return m.GetButtonsMessage().GetContextInfo().StanzaId
	} else if m.GetGroupInviteMessage().GetContextInfo() != nil {
		return m.GetGroupInviteMessage().GetContextInfo().StanzaId
	} else if m.GetProductMessage().GetContextInfo() != nil {
		return m.GetProductMessage().GetContextInfo().StanzaId
	} else if m.GetListMessage().GetContextInfo() != nil {
		return m.GetListMessage().GetContextInfo().StanzaId
	} else if m.GetTemplateMessage().GetContextInfo() != nil {
		return m.GetTemplateMessage().GetContextInfo().StanzaId
	} else if m.GetContactMessage().GetContextInfo() != nil {
		return m.GetContactMessage().GetContextInfo().StanzaId
	} else {
		return &evm.Info.ID
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
	case "location":
		return m.GetAudioMessage().GetContextInfo().GetQuotedMessage()
	default:
		return ParseQuotedMessage(m)
	}
}

func WithReply(m *events.Message) *waProto.ContextInfo {
	return &waProto.ContextInfo{
		StanzaId:      &m.Info.ID,
		Participant:   types.String(m.Info.MessageSource.Sender.String()),
		QuotedMessage: m.Message,
	}
}

func ParseJID(arg string) (waTypes.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return waTypes.NewJID(arg, waTypes.DefaultUserServer), true
	} else {
		recipient, err := waTypes.ParseJID(arg)
		if err != nil {
			return recipient, false
		} else if recipient.User == "" {
			return recipient, false
		}
		return recipient, true
	}
}

func RevokeMessage(c *whatsmeow.Client, jid waTypes.JID, messageId waTypes.MessageID) {
	_, err := c.SendMessage(context.Background(), jid, c.BuildRevoke(jid, waTypes.EmptyJID, messageId))
	if err != nil {
		fmt.Printf("error: revoking message: %v\n", err)
	}
	return
}
