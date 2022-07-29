package handler

import (
	"github.com/itzngga/goRoxy/helper"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
	"regexp"
	"strings"
)

var cmdRegex = regexp.MustCompile(`!|#|\?|.|@|&|\*|-|=|\+`)

func ParseCmd(str string) (string, bool) {
	if str == "" {
		return "", false
	}
	word := strings.Fields(str)
	if len(word) == 0 && word[0] == "" {
		return "", false
	}
	if cmdRegex.MatchString(word[0][0:1]) {
		return word[0][1:], true
	}
	return "", false
}

func ParseMessageText(m *events.Message) string {
	var pesan *waProto.Message
	if m.IsViewOnce {
		pesan = m.Message.GetViewOnceMessage().GetMessage()
	} else if m.IsEphemeral {
		pesan = m.Message.GetEphemeralMessage().GetMessage()
	} else {
		pesan = m.Message
	}
	if pesan.GetConversation() != "" {
		return pesan.GetConversation()
	} else if pesan.GetVideoMessage().GetCaption() != "" {
		return pesan.GetVideoMessage().GetCaption()
	} else if pesan.GetImageMessage().GetCaption() != "" {
		return pesan.GetImageMessage().GetCaption()
	} else if pesan.GetExtendedTextMessage().GetText() != "" {
		return pesan.GetExtendedTextMessage().GetText()
	} else if pesan.GetTemplateButtonReplyMessage().GetSelectedId() != "" {
		cmd, err := helper.XTCrypto.MakeDecrypt(pesan.GetTemplateButtonReplyMessage().GetSelectedId())
		if err != nil {
			return ""
		}
		return cmd
	} else if pesan.GetButtonsResponseMessage().GetSelectedButtonId() != "" {
		cmd, err := helper.XTCrypto.MakeDecrypt(pesan.GetButtonsResponseMessage().GetSelectedButtonId())
		if err != nil {
			return ""
		}
		return cmd
	} else if pesan.GetListResponseMessage().GetSingleSelectReply().GetSelectedRowId() != "" {
		cmd, err := helper.XTCrypto.MakeDecrypt(pesan.GetListResponseMessage().GetSingleSelectReply().GetSelectedRowId())
		if err != nil {
			return ""
		}
		return cmd
	} else {
		return ""
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

func RemoveElementByIndex[T any](slice []T, index int) []T {
	sliceLen := len(slice)
	sliceLastIndex := sliceLen - 1
	if index != sliceLastIndex {
		slice[index] = slice[sliceLastIndex]
	}
	return slice[:sliceLastIndex]
}

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
