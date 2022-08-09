package middleware

import (
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/jellydator/ttlcache/v2"
	"os"
	"strconv"
	"time"
)

func GenerateAllMiddlewares() {
	PrepareMiddleware()
	AddMiddleware(LogMiddleware)
	AddMiddleware(CooldownMiddleware)
}
func PrepareMiddleware() {
	cd, _ := strconv.Atoi(os.Getenv("DEFAULT_COOLDOWN_SEC"))
	cooldownCache = ttlcache.NewCache()
	cooldownTimeout = time.Duration(cd) * time.Second
}

func AddMiddleware(mid handler.MiddlewareFunc) {
	handler.GlobalMiddleware = append(handler.GlobalMiddleware, mid)
}
