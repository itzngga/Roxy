package command

import (
	"github.com/itzngga/goRoxy/options"
	"github.com/zhangyunhao116/skipmap"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"time"
)

type RunFuncParams struct {
	ParsedMsg string
	Arguments []string
	Number    string

	Options *options.Options
	Log     waLog.Logger
	Event   *events.Message
	Info    *waTypes.MessageInfo
	User    *waTypes.JID
	Message *waProto.Message
	Cmd     *Command

	Locals *skipmap.StringMap[string]
}

func (r *RunFuncParams) GetLocals(key string) (string, bool) {
	return r.Locals.Load(key)
}

func (r *RunFuncParams) SetLocals(key string, value string) {
	r.Locals.Store(key, value)
	return
}

func (r *RunFuncParams) SetLocalsWithTTL(key string, value string, ttl time.Duration) {
	r.Locals.Store(key, value)
	go func() {
		timeout := time.NewTimer(ttl)
		<-timeout.C
		r.Locals.Delete(key)
		timeout.Stop()
	}()
	return
}

func (r *RunFuncParams) Del(key string) {
	r.Locals.Delete(key)
	return
}
