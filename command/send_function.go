package command

import (
	"github.com/itzngga/Roxy/types"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func (runFunc *RunFuncContext) GetAllChats() []*events.Message {
	GetAllChats := types.GetContext[types.GetAllChats](runFunc.Ctx, "GetAllChats")
	return GetAllChats()
}

func (runFunc *RunFuncContext) GetChatInJID(jid waTypes.JID) []*events.Message {
	GetChatInJID := types.GetContext[types.GetChatInJID](runFunc.Ctx, "GetChatInJID")
	return GetChatInJID(jid)
}

func (runFunc *RunFuncContext) GetStatusMessages() []*events.Message {
	GetStatusMessages := types.GetContext[types.GetStatusMessages](runFunc.Ctx, "GetStatusMessages")
	return GetStatusMessages()
}

func (runFunc *RunFuncContext) FindMessageByID(jid waTypes.JID, id string) *events.Message {
	FindMessageByID := types.GetContext[types.FindMessageByID](runFunc.Ctx, "FindMessageByID")
	return FindMessageByID(jid, id)
}
