package middleware

import (
	"github.com/itzngga/goRoxy/command"
	"github.com/zhangyunhao116/skipmap"
	"go.mau.fi/whatsmeow"
	"time"
)

var CooldownCache *skipmap.StringMap[string]

func CooldownMiddleware(c *whatsmeow.Client, args command.RunFuncArgs) bool {
	if args.Options.CommandCooldownTimeout == 0 {
		return true
	}
	_, ok := CooldownCache.Load(c.Store.ID.User + "-" + args.Evm.Info.Sender.User)
	if ok {
		//util.SendReplyMessage(c, args.Evm, "You are on Cooldown!")
		return false
	}
	go func() {
		timeout := time.NewTimer(args.Options.CommandCooldownTimeout)
		<-timeout.C
		CooldownCache.Delete(c.Store.ID.User + "-" + args.Evm.Info.Sender.User)
		timeout.Stop()
	}()

	return true
}
