package command

import (
	"fmt"
	waTypes "go.mau.fi/whatsmeow/types"
	"time"
)

// SendReadPresence send read status in current chat
func (runFunc *RunFuncContext) SendReadPresence() {
	jids := []waTypes.MessageID{
		runFunc.MessageInfo.ID,
	}
	runFunc.Client.MarkRead(jids, runFunc.MessageInfo.Timestamp, runFunc.MessageInfo.Chat, runFunc.MessageInfo.Sender)
	return
}

// SendTypingPresence send typing action in current chat
func (runFunc *RunFuncContext) SendTypingPresence(duration time.Duration) {
	go func() {
		chat := runFunc.MessageInfo.Chat
		err := runFunc.Client.SubscribePresence(chat)
		if err != nil {
			fmt.Println(err)
		}

		err = runFunc.Client.SendChatPresence(chat, "composing", "")
		if err != nil {
			fmt.Println(err)
		}

		if duration != 0 {
			time.Sleep(duration)
		} else {
			time.Sleep(1500 * time.Millisecond)
		}

		err = runFunc.Client.SendChatPresence(chat, "paused", "")
		if err != nil {
			fmt.Println(err)
		}
	}()

}
