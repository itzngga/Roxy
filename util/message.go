package util

import (
	"strings"

	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
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

func ParseQuotedRemoteJid(m *waProto.Message) *string {
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
		return nil
	}
}

func ParseMentionedJid(m *waProto.Message) []string {
	if m.GetExtendedTextMessage().GetContextInfo() != nil {
		return m.GetExtendedTextMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetImageMessage().GetContextInfo() != nil {
		return m.GetImageMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetVideoMessage().GetContextInfo() != nil {
		return m.GetVideoMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetDocumentMessage().GetContextInfo() != nil {
		return m.GetDocumentMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetAudioMessage().GetContextInfo() != nil {
		return m.GetAudioMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetStickerMessage().GetContextInfo() != nil {
		return m.GetStickerMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetButtonsMessage().GetContextInfo() != nil {
		return m.GetButtonsMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetGroupInviteMessage().GetContextInfo() != nil {
		return m.GetGroupInviteMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetProductMessage().GetContextInfo() != nil {
		return m.GetProductMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetListMessage().GetContextInfo() != nil {
		return m.GetListMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetTemplateMessage().GetContextInfo() != nil {
		return m.GetTemplateMessage().GetContextInfo().GetMentionedJID()
	} else if m.GetContactMessage().GetContextInfo() != nil {
		return m.GetContactMessage().GetContextInfo().GetMentionedJID()
	} else {
		return make([]string, 0)
	}
}

func ParseQuotedMessageId(m *waProto.Message) *string {
	if m.GetExtendedTextMessage().GetContextInfo() != nil {
		return m.GetExtendedTextMessage().GetContextInfo().StanzaID
	} else if m.GetImageMessage().GetContextInfo() != nil {
		return m.GetImageMessage().GetContextInfo().StanzaID
	} else if m.GetVideoMessage().GetContextInfo() != nil {
		return m.GetVideoMessage().GetContextInfo().StanzaID
	} else if m.GetDocumentMessage().GetContextInfo() != nil {
		return m.GetDocumentMessage().GetContextInfo().StanzaID
	} else if m.GetAudioMessage().GetContextInfo() != nil {
		return m.GetAudioMessage().GetContextInfo().StanzaID
	} else if m.GetStickerMessage().GetContextInfo() != nil {
		return m.GetStickerMessage().GetContextInfo().StanzaID
	} else if m.GetButtonsMessage().GetContextInfo() != nil {
		return m.GetButtonsMessage().GetContextInfo().StanzaID
	} else if m.GetGroupInviteMessage().GetContextInfo() != nil {
		return m.GetGroupInviteMessage().GetContextInfo().StanzaID
	} else if m.GetProductMessage().GetContextInfo() != nil {
		return m.GetProductMessage().GetContextInfo().StanzaID
	} else if m.GetListMessage().GetContextInfo() != nil {
		return m.GetListMessage().GetContextInfo().StanzaID
	} else if m.GetTemplateMessage().GetContextInfo() != nil {
		return m.GetTemplateMessage().GetContextInfo().StanzaID
	} else if m.GetContactMessage().GetContextInfo() != nil {
		return m.GetContactMessage().GetContextInfo().StanzaID
	} else {
		return nil
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
		StanzaID:      &m.Info.ID,
		Participant:   proto.String(m.Info.Sender.String()),
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
