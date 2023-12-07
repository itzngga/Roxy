package types

import (
	"github.com/go-whatsapp/whatsmeow"
	waProto "github.com/go-whatsapp/whatsmeow/binary/proto"
	waTypes "github.com/go-whatsapp/whatsmeow/types"
	"github.com/go-whatsapp/whatsmeow/types/events"
)

type AppMethods interface {
	UpsertMessages(jid waTypes.JID, message []*events.Message)
	GetAllChats() []*events.Message
	GetChatInJID(jid waTypes.JID) []*events.Message
	GetStatusMessages() []*events.Message
	FindMessageByID(jid waTypes.JID, id string) *events.Message
	SendMessage(to waTypes.JID, message *waProto.Message, extra ...whatsmeow.SendRequestExtra) (whatsmeow.SendResponse, error)
	ClientJID() waTypes.JID
	Client() *whatsmeow.Client
}

type MuxerMethods interface {
	FindGroupByJid(groupJid waTypes.JID) (group *waTypes.GroupInfo, err error)
	GetAllGroups() (group []*waTypes.GroupInfo, err error)
	UnCacheOneGroup(info *events.GroupInfo, joined *events.JoinedGroup)
	IsGroupAdmin(chat waTypes.JID, jid any) (bool, error)
	IsClientGroupAdmin(chat waTypes.JID) (bool, error)
	SendEmojiMessage(event *events.Message, emoji string)
	CacheAllGroup()
}

type ClientMethods interface {
	AppMethods
	MuxerMethods
}
