package message_middleware

import (
	"github.com/itzngga/goRoxy/command"
	"github.com/zhangyunhao116/skipmap"
	"go.mau.fi/whatsmeow"
	"time"
)

var CooldownCache *skipmap.StringMap[string]

func CooldownMiddleware(c *whatsmeow.Client, params *command.RunFuncParams) bool {
	if params.Options.CommandCooldownTimeout == 0 {
		return true
	}
	_, ok := CooldownCache.Load(c.Store.ID.User + "-" + params.Event.Info.Sender.User)
	if ok {
		//util.SendReplyMessage(c, params.Event, "You are on Cooldown!")
		return false
	}
	go func() {
		timeout := time.NewTimer(params.Options.CommandCooldownTimeout)
		<-timeout.C
		CooldownCache.Delete(c.Store.ID.User + "-" + params.Event.Info.Sender.User)
		timeout.Stop()
	}()

	return true
}
