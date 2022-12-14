package command

import (
	"fmt"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"sort"
)

type MiddlewareFunc func(c *whatsmeow.Client, args RunFuncArgs) bool
type RunFunc func(c *whatsmeow.Client, args RunFuncArgs) *waProto.Message

type Command struct {
	Name        string
	Aliases     []string
	Description string

	Category string
	Cache    bool

	HideFromHelp bool
	GroupOnly    bool
	PrivateOnly  bool
	Middleware   MiddlewareFunc
	RunFunc      RunFunc
}

func (c *Command) GetName(name string) string {
	var theName string
	if c.Name == name {
		theName = name
	}
	if cmd := sort.SearchStrings(c.Aliases, name); cmd != len(c.Aliases) {
		theName = c.Aliases[cmd]
	}
	return theName
}

func (c *Command) Validate() {
	if c.Name == "" {
		panic("Command name cannot be empty")
	} else if c.Description == "" {
		c.Description = fmt.Sprintf("This is %s command description example", c.Name)
	} else if c.RunFunc == nil {
		panic("RunFunc cannot be empty")
	}

	sort.Strings(c.Aliases)
}
