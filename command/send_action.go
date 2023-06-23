package command

import waTypes "go.mau.fi/whatsmeow/types"

func (runFunc *RunFuncContext) SendReadPresence() {
	jids := []waTypes.MessageID{
		runFunc.MessageInfo.ID,
	}
	runFunc.Client.MarkRead(jids, runFunc.MessageInfo.Timestamp, runFunc.MessageInfo.Chat, runFunc.MessageInfo.Sender)
	return
}
