package util

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"regexp"
	"strings"
)

var cmdRegex = regexp.MustCompile(`!|#|\?|.|@|&|\*|-|=|\+`)

func ParseCmd(str string) (string, bool) {
	if str == "" {
		return "", false
	}
	word := strings.Split(str, " ")
	if len(word) == 0 && word[0] == "" {
		return "", false
	}
	if cmdRegex.MatchString(word[0][0:1]) {
		return word[0][1:], true
	}
	return "", false
}

func ParseMessageText(uid string, m *events.Message) string {
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

		// Button,Template,List Was Deprecated
		//} else if pesan.GetTemplateButtonReplyMessage().GetSelectedId() != "" {
		//	cmd := ParseButtonID(uid, pesan.GetTemplateButtonReplyMessage().GetSelectedId())
		//	return cmd
		//} else if pesan.GetButtonsResponseMessage().GetSelectedButtonId() != "" {
		//	cmd := ParseButtonID(uid, pesan.GetButtonsResponseMessage().GetSelectedButtonId())
		//	return cmd
		//} else if pesan.GetListResponseMessage().GetSingleSelectReply().GetSelectedRowId() != "" {
		//	cmd := ParseButtonID(uid, pesan.GetListResponseMessage().GetSingleSelectReply().GetSelectedRowId())
		//	return cmd
		//} else if pesan.GetTemplateButtonReplyMessage().GetSelectedId() != "" {
		//	cmd := ParseButtonID(uid, pesan.GetTemplateButtonReplyMessage().GetSelectedId())
		//	return cmd
	} else {
		return ""
	}
}
