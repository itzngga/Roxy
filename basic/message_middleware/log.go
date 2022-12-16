package message_middleware

import (
	"github.com/itzngga/goRoxy/command"
	"go.mau.fi/whatsmeow"
)

func LogMiddleware(c *whatsmeow.Client, params *command.RunFuncParams) bool {
	params.Log.Infof("[CMD] [%s] command > %s", params.Number, params.Cmd.Name)
	return true
}
