package command

import (
	"fmt"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"sort"
	"time"
)

type MiddlewareFunc func(c *RunFuncContext) bool
type RunFunc func(c *RunFuncContext) *waProto.Message
type StateFunc func(c *RunFuncContext, parsed string)

type StateCommand struct {
	Name        string
	GroupOnly   bool
	PrivateOnly bool
	CancelReply string

	StateTimeout time.Duration
	RunFunc      StateFunc
}

func (c *StateCommand) Validate() {
	if c.Name == "" {
		panic("error: command name cannot be empty")
	}
	if c.RunFunc == nil {
		panic("error: RunFunc cannot be empty")
	}
	if c.PrivateOnly && c.GroupOnly {
		panic("error: invalid scope group/private?")
	}
	if c.CancelReply == "" {
		c.CancelReply = "User cancelled the State"
	}
	if c.StateTimeout == 0 {
		c.StateTimeout = time.Minute * 15
	}
}

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
