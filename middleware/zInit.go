package middleware

import "github.com/itzngga/goRoxy/internal/handler"

func GenerateAllMiddlewares() {
	AddMiddleware(LogMiddleware)
}

func AddMiddleware(mid handler.MiddlewareFunc) {
	handler.GlobalMiddleware = append(handler.GlobalMiddleware, mid)
}
