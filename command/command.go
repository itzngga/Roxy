package command

import (
	"fmt"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"sort"
)

type MiddlewareFunc func(c *RunFuncContext) bool
type RunFunc func(c *RunFuncContext) *waProto.Message
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

	Middleware MiddlewareFunc
	RunFunc    RunFunc
}

func (c *Command) Validate() {
	if c.Name == "" {
		panic("error: command name cannot be empty")
	}
	if c.Description == "" {
		c.Description = fmt.Sprintf("This is %s command description example", c.Name)
	}
	if c.RunFunc == nil {
		panic("error: RunFunc cannot be empty")
	}
	if c.PrivateOnly && c.GroupOnly {
		panic("error: invalid scope group/private?")
	}

	sort.Strings(c.Aliases)
}
