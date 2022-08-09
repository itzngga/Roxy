package middleware

import (
	"fmt"
	"github.com/itzngga/goRoxy/internal/handler"
	"go.mau.fi/whatsmeow"
)

func LogMiddleware(c *whatsmeow.Client, args handler.RunFuncArgs) bool {
	fmt.Println("[CMD] Command : " + args.Cmd.Name)
	return true
}
