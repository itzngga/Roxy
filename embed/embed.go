package embed

import (
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/types"
)

var Commands types.Embed[*command.Command]
var GlobalMiddlewares types.Embed[command.MiddlewareFunc]
var Middlewares types.Embed[command.MiddlewareFunc]
var Categories types.Embed[string]

func init() {
	cmd := types.NewEmbed[*command.Command]()
	mid := types.NewEmbed[command.MiddlewareFunc]()
	cat := types.NewEmbed[string]()
	gMid := types.NewEmbed[command.MiddlewareFunc]()

	Commands = &cmd
	Middlewares = &mid
	Categories = &cat
	GlobalMiddlewares = &gMid
}
