package context

import (
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// GetAllChats get all chats from local store
func (context *Ctx) GetAllChats() []*events.Message {
	return context.Methods().GetAllChats()
}

// GetChatInJID get chats in jid from local store
func (context *Ctx) GetChatInJID(jid waTypes.JID) []*events.Message {
	return context.Methods().GetChatInJID(jid)
}

// GetStatusMessages get client statuses from local store
func (context *Ctx) GetStatusMessages() []*events.Message {
	return context.Methods().GetStatusMessages()
}

// FindMessageByID find message by specific id from local store
func (context *Ctx) FindMessageByID(jid waTypes.JID, id string) *events.Message {
	return context.Methods().FindMessageByID(jid, id)
}
