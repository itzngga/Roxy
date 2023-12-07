package roxy

import (
	"fmt"
	"github.com/itzngga/Roxy/context"
	"github.com/itzngga/Roxy/types"
	"sort"
)

var Commands types.Embed[*Command]
var GlobalMiddlewares types.Embed[context.MiddlewareFunc]
var Middlewares types.Embed[context.MiddlewareFunc]
var Categories types.Embed[string]

func init() {
	cmd := types.NewEmbed[*Command]()
	mid := types.NewEmbed[context.MiddlewareFunc]()
	cat := types.NewEmbed[string]()
	gMid := types.NewEmbed[context.MiddlewareFunc]()

	Commands = &cmd
	Middlewares = &mid
	Categories = &cat
	GlobalMiddlewares = &gMid
}

type Command struct {
	Name        string
	Aliases     []string
	Description string

	Category string
	Cache    bool

	HideFromHelp bool
	GroupOnly    bool
	PrivateOnly  bool

	OnlyAdminGroup   bool
	OnlyIfBotAdmin   bool
	AdditionalValues map[string]interface{}

	RunFunc    context.RunFunc
	Middleware context.MiddlewareFunc
}

func (c *Command) Validate() {
	if c.Name == "" {
		panic("error: command name cannot be empty")
	}
	if c.Description == "" {
		c.Description = fmt.Sprintf("This is %s command description example", c.Name)
	}
	if c.PrivateOnly && c.GroupOnly {
		panic("error: invalid scope group/private?")
	}

	sort.Strings(c.Aliases)
}
