package command

import (
	"fmt"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"sort"
)

type MiddlewareFunc func(c *whatsmeow.Client, params *RunFuncParams) bool
type RunFunc func(c *whatsmeow.Client, params *RunFuncParams) *waProto.Message

type Command struct {
	Name        string
	Aliases     []string
	Description string

	Category string
	Cache    bool
	BuiltIn  bool

	HideFromHelp bool
	GroupOnly    bool
	PrivateOnly  bool
	Middleware   MiddlewareFunc
	RunFunc      RunFunc
}

func (c *Command) Validate() {
	if c.Name == "" {
		panic("error: command name cannot be empty")
	} else if c.Description == "" {
		c.Description = fmt.Sprintf("This is %s command description example", c.Name)
	} else if c.RunFunc == nil {
		panic("error: RunFunc cannot be empty")
	}
	if c.PrivateOnly && c.GroupOnly {
		panic("error: invalid scope group/private?")
	}

	sort.Strings(c.Aliases)
}
