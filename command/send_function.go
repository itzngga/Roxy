package command

import (
	"github.com/itzngga/Roxy/types"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// GetAllChats get all chats from local store
func (runFunc *RunFuncContext) GetAllChats() []*events.Message {
	GetAllChats := types.GetContext[types.GetAllChats](runFunc.Ctx, "GetAllChats")
	return GetAllChats()
}

// GetChatInJID get chats in jid from local store
func (runFunc *RunFuncContext) GetChatInJID(jid waTypes.JID) []*events.Message {
	GetChatInJID := types.GetContext[types.GetChatInJID](runFunc.Ctx, "GetChatInJID")
	return GetChatInJID(jid)
}

// GetStatusMessages get client statuses from local store
func (runFunc *RunFuncContext) GetStatusMessages() []*events.Message {
	GetStatusMessages := types.GetContext[types.GetStatusMessages](runFunc.Ctx, "GetStatusMessages")
	return GetStatusMessages()
}

// FindMessageByID find message by specific id from local store
func (runFunc *RunFuncContext) FindMessageByID(jid waTypes.JID, id string) *events.Message {
	FindMessageByID := types.GetContext[types.FindMessageByID](runFunc.Ctx, "FindMessageByID")
	return FindMessageByID(jid, id)
}
