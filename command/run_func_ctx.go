package command

import (
	"github.com/itzngga/Roxy/options"
	"github.com/itzngga/Roxy/types"
	"github.com/zhangyunhao116/skipmap"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"time"
)

type RunFuncContext struct {
	FromMe    bool
	ParsedMsg string
	Number    string
	Prefix    string
	Arguments []string

	CurrentCommand *Command
	WaLog          waLog.Logger
	ClientJID      *waTypes.JID
	MessageEvent   *events.Message
	Options        *options.Options
	Message        *waProto.Message
	Client         *whatsmeow.Client
	QuestionChan   chan *QuestionState
	MessageInfo    *waTypes.MessageInfo
	Locals         *skipmap.StringMap[string]
	Ctx            *skipmap.StringMap[types.RoxyContext]
}

func (runFunc *RunFuncContext) GetClient() *whatsmeow.Client {
	return runFunc.Client
}

func (runFunc *RunFuncContext) GetOptions() *options.Options {
	return runFunc.Options
}

func (runFunc *RunFuncContext) GetMessageEvent() *events.Message {
	return runFunc.MessageEvent
}

func (runFunc *RunFuncContext) GetMessageInfo() *waTypes.MessageInfo {
	return runFunc.MessageInfo
}

func (runFunc *RunFuncContext) GetClientJID() *waTypes.JID {
	return runFunc.ClientJID
}

func (runFunc *RunFuncContext) GetCommand() *Command {
	return runFunc.CurrentCommand
}
func (runFunc *RunFuncContext) GetMessage() *waProto.Message {
	return runFunc.Message
}

func (runFunc *RunFuncContext) RangeLocals(fun func(key string, value string) bool) {
	runFunc.Locals.Range(fun)
}

func (runFunc *RunFuncContext) GetLocals(key string) (string, bool) {
	return runFunc.Locals.Load(key)
}

func (runFunc *RunFuncContext) SetLocals(key string, value string) {
	runFunc.Locals.Store(key, value)
	return
}

func (runFunc *RunFuncContext) DelLocals(key string) {
	runFunc.Locals.Delete(key)
	return
}

func (runFunc *RunFuncContext) SetLocalsWithTTL(key string, value string, ttl time.Duration) {
	runFunc.Locals.Store(key, value)
	go func() {
		timeout := time.NewTimer(ttl)
		<-timeout.C
		runFunc.Locals.Delete(key)
		timeout.Stop()
	}()
	return
}
