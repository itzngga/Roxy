package middleware

import (
	"fmt"
	"github.com/itzngga/goRoxy/command"
	"go.mau.fi/whatsmeow"
)

func LogMiddleware(c *whatsmeow.Client, args command.RunFuncArgs) bool {
	fmt.Println(fmt.Sprintf("[CMD] [%s] command : %s", args.Number, args.Cmd.Name))
	return true
}
