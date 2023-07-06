package util

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"regexp"
	"strings"
)

var cmdRegex = regexp.MustCompile("^(!|#|\\?|.|@|&|\\*|-|=|\\+|/)?(.*)")

var ListPrefixes = []string{"!", "/", "#", "?", ".", "@", "&", "*", "-", "=", "+"}

func ParseCmd(str string) (prefix string, cmd string, ok bool) {
	if str == "" {
		return "", "", false
	}
	split := strings.Split(str, " ")
	switch string(str[0]) {
	case "!":
		return string(str[0]), split[0][1:], true
	case "/":
		return string(str[0]), split[0][1:], true
	case "#":
		return string(str[0]), split[0][1:], true
	case "?":
		return string(str[0]), split[0][1:], true
	case ".":
		return string(str[0]), split[0][1:], true
	case "@":
		return string(str[0]), split[0][1:], true
	case "&":
		return string(str[0]), split[0][1:], true
	case "*":
		return string(str[0]), split[0][1:], true
	case "-":
		return string(str[0]), split[0][1:], true
	case "=":
		return string(str[0]), split[0][1:], true
	case "+":
		return string(str[0]), split[0][1:], true
	}

	if len(str) >= 1 {
		return string(str[0]), split[0][1:], false
	}
	//word := strings.Split(str, " ")
	//if len(word) == 0 && word[0] == "" {
	//	return "", "", false
	//}
	//
	//if len(word) >= 1 && word[0] != "" && cmdRegex.MatchString(string(word[0][0])) {
	//	fmt.Println(cmdRegex.FindAllString(string(word[0]), -1))
	//	return string(word[0][0]), word[0][1:], true
	//}
	//if cmdRegex.MatchString(word[0][0:1]) {
	//	return string(word[0][0]), word[0][1:], true
	//}
	return "", "", false
}

func GetQuotedText(m *events.Message) string {
	var pesan *waProto.Message
	if m.IsViewOnce {
		pesan = m.Message.GetViewOnceMessage().GetMessage()
	} else if m.IsEphemeral {
		pesan = m.Message.GetEphemeralMessage().GetMessage()
	} else {
		pesan = m.Message
	}

	if quoted := ParseQuotedBy(pesan, "text"); quoted != nil {
		if quoted.GetExtendedTextMessage() != nil {
			return *quoted.ExtendedTextMessage.Text
		} else {
			return *quoted.Conversation
		}
	} else {
		if pesan.GetExtendedTextMessage() != nil {
			return *pesan.ExtendedTextMessage.Text
		} else {
			return *pesan.Conversation
		}
	}
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
