package message_middleware

import (
	"fmt"
	"github.com/itzngga/goRoxy/command"
	"go.mau.fi/whatsmeow"
	"time"
)

func LogMiddleware(c *whatsmeow.Client, args command.RunFuncArgs) bool {
	fmt.Println(fmt.Sprintf("[%s] [CMD] [%s] command : %s", time.Now().Format("03:04:05"), args.Number, args.Cmd.Name))
	return true
}
