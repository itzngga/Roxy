package middleware

import (
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"github.com/jellydator/ttlcache/v2"
	"go.mau.fi/whatsmeow"
	"time"
)

var cooldownCache ttlcache.SimpleCache
var cooldownTimeout time.Duration

func CooldownMiddleware(c *whatsmeow.Client, args handler.RunFuncArgs) bool {
	if cooldownTimeout == 0 {
		return true
	}
	cdId, _ := cooldownCache.Get(c.Store.ID.User + args.Evm.Info.Sender.User)
	if cdId != nil {
		util.SendReplyMessage(c, args.Evm, "You are on Cooldown!")
		return false
	}
	go func() {
		err := cooldownCache.SetWithTTL(c.Store.ID.User+args.Evm.Info.Sender.User, true, cooldownTimeout)
		if err != nil {
			return
		}
	}()

	return true
}
