package embed

import (
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/types"
)

var Commands types.Embed[*command.Command]
var GlobalMiddlewares types.Embed[command.MiddlewareFunc]
var Middlewares types.Embed[command.MiddlewareFunc]
var Categories types.Embed[string]

func init() {
	cmd := types.NewEmbed[*command.Command]()
	mid := types.NewEmbed[command.MiddlewareFunc]()
	cat := types.NewEmbed[string]()
	gmid := types.NewEmbed[command.MiddlewareFunc]()

	Commands = &cmd
	Middlewares = &mid
	Categories = &cat
	GlobalMiddlewares = &gmid
}
