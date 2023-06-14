package command

import (
	"github.com/itzngga/Roxy/options"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type StateFuncContext struct {
	RunFuncContext
	ParsedMsg string
	Number    string

	Client       *whatsmeow.Client
	Options      *options.Options
	MessageEvent *events.Message
	MessageInfo  *waTypes.MessageInfo
	ClientJID    *waTypes.JID
	Message      *waProto.Message
	CurrentState *StateCommand
	WaLog        waLog.Logger

	Locals map[string]interface{}
}

func (runFunc *StateFuncContext) GetClient() *whatsmeow.Client {
	return runFunc.Client
}

func (runFunc *StateFuncContext) GetOptions() *options.Options {
	return runFunc.Options
}

func (runFunc *StateFuncContext) GetMessageEvent() *events.Message {
	return runFunc.MessageEvent
}

func (runFunc *StateFuncContext) GetMessageInfo() *waTypes.MessageInfo {
	return runFunc.MessageInfo
}

func (runFunc *StateFuncContext) GetClientJID() *waTypes.JID {
	return runFunc.ClientJID
}

func (runFunc *StateFuncContext) GetState() *StateCommand {
	return runFunc.CurrentState
}
func (runFunc *StateFuncContext) GetMessage() *waProto.Message {
	return runFunc.Message
}
