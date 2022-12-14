package command

import (
	"github.com/itzngga/goRoxy/options"
	"github.com/zhangyunhao116/skipmap"
	"go.mau.fi/whatsmeow/types/events"
	"time"
)

type RunFuncArgs struct {
	Options *options.Options
	Evm     *events.Message
	Cmd     *Command
	Msg     string
	Args    []string
	Number  string
	Locals  *skipmap.StringMap[string]
}

func (r *RunFuncArgs) GetLocals(key string) (string, bool) {
	return r.Locals.Load(key)
}

func (r *RunFuncArgs) SetLocals(key string, value string) {
	r.Locals.Store(key, value)
	return
}

func (r *RunFuncArgs) SetLocalsWithTTL(key string, value string, ttl time.Duration) {
	r.Locals.Store(key, value)
	go func() {
		timeout := time.NewTimer(ttl)
		<-timeout.C
		r.Locals.Delete(key)
		timeout.Stop()
	}()
	return
}

func (r *RunFuncArgs) Del(key string) {
	r.Locals.Delete(key)
	return
}
