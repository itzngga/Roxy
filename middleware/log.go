package middleware

import (
	"github.com/itzngga/goRoxy/internal/handler"
	"go.mau.fi/whatsmeow"
	"go.uber.org/zap"
)

func LogMiddleware(c *whatsmeow.Client, args handler.RunFuncArgs) bool {
	//fmt.Println(fmt.Sprintf("[CMD] [%s] command : %s", args.Number, args.Cmd.Name))
	go cmdLogger.Info("command", zap.String("number", args.Number), zap.String("cmd", args.Cmd.Name))
	return true
}
