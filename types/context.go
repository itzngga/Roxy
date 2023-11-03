package types

import (
	"github.com/zhangyunhao116/skipmap"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type RoxyContext interface{}

// App Context type
type UpsertMessages func(jid waTypes.JID, message []*events.Message)
type GetAllChats func() []*events.Message
type GetChatInJID func(jid waTypes.JID) []*events.Message
type GetStatusMessages func() []*events.Message
type FindMessageByID func(jid waTypes.JID, id string) *events.Message
type SendMessage func(to waTypes.JID, message *waProto.Message, extra ...whatsmeow.SendRequestExtra) (whatsmeow.SendResponse, error)

// muxer Context type
type FindGroupByJid func(groupJid waTypes.JID) (group *waTypes.GroupInfo, err error)
type GetAllGroups func() (group []*waTypes.GroupInfo, err error)
type CacheAllGroup func()
type UNCacheOneGroup func(info *events.GroupInfo, joined *events.JoinedGroup)
type IsGroupAdmin func(chat waTypes.JID, jid any) (bool, error)
type IsClientGroupAdmin func(chat waTypes.JID) (bool, error)
type SendEmojiMessage func(event *events.Message, emoji string)

// GetContext dynamically get context value
func GetContext[T any](ctx *skipmap.StringMap[RoxyContext], key string) T {
	load, _ := ctx.Load(key)
	return load.(T)
}
