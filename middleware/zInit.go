package middleware

import (
	"github.com/itzngga/goRoxy/config"
	"github.com/itzngga/goRoxy/helper"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/jellydator/ttlcache/v2"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

var cmdLogger *zap.Logger

func GenerateAllMiddlewares() {
	PrepareMiddleware()
	AddMiddleware(LogMiddleware)
	AddMiddleware(CooldownMiddleware)
}
func PrepareMiddleware() {
	cd, _ := strconv.Atoi(os.Getenv("DEFAULT_COOLDOWN_SEC"))
	cooldownCache = ttlcache.NewCache()
	cooldownTimeout = time.Duration(cd) * time.Second
	cmdLogger = config.CmdLogger("info")
}

func AddMiddleware(mid handler.MiddlewareFunc) {
	handler.GlobalMiddleware.Store(helper.CreateUid(), mid)
}
