package context

import (
	waTypes "go.mau.fi/whatsmeow/types"
	"time"
)

// SendReadPresence send read status in current chat
func (context *Ctx) SendReadPresence() {
	jids := []waTypes.MessageID{
		context.info.ID,
	}
	context.client.MarkRead(jids, context.info.Timestamp, context.info.Chat, context.info.Sender)
	return
}

// SendTypingPresence send typing action in current chat
func (context *Ctx) SendTypingPresence(duration time.Duration) {
	go func() {
		chat := context.info.Chat
		err := context.client.SubscribePresence(chat)
		if err != nil {
			context.logger.Errorf("error: v", err)
			return
		}

		err = context.client.SendChatPresence(chat, "composing", "")
		if err != nil {
			context.logger.Errorf("error: v", err)
			return
		}

		if duration != 0 {
			time.Sleep(duration)
		} else {
			time.Sleep(1500 * time.Millisecond)
		}

		err = context.client.SendChatPresence(chat, "paused", "")
		if err != nil {
			context.logger.Errorf("error: v", err)
			return
		}
	}()
}
