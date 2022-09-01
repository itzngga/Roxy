package handler

import (
	"fmt"
	"github.com/zhangyunhao116/skipmap"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"sort"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type MiddlewareFunc func(c *whatsmeow.Client, args RunFuncArgs) bool
type RunFunc func(c *whatsmeow.Client, args RunFuncArgs) *waProto.Message

type RunFuncArgs struct {
	Evm    *events.Message
	Cmd    *Command
	Msg    string
	Args   []string
	Number string
	Locals *skipmap.StringMap[string]
}

func (r *RunFuncArgs) GetLocals(key string) (string, bool) {
	return r.Locals.Load(key)
}

func (r *RunFuncArgs) SetLocals(key string, value string) {
	r.Locals.Store(key, value)
	return
}

func (r *RunFuncArgs) SetLocalsWithTTL(key string, value string, ttl time.Duration) {
	r.Locals.Store(key, value)
	go func() {
		timeout := time.NewTimer(ttl)
		<-timeout.C
		r.Locals.Delete(key)
		timeout.Stop()
	}()
	return
}

func (r *RunFuncArgs) Del(key string) {
	r.Locals.Delete(key)
	return
}

type Command struct {
	Name            string
	Aliases         []string
	Description     string
	LongDescription string

	Cooldown time.Duration
	Category *Category

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
		c.Description = fmt.Sprintf("This is %s command description", c.Name)
	} else if c.LongDescription == "" {
		c.LongDescription = fmt.Sprintf("This is %s command long description", c.Name)
	} else if c.Cooldown == 0 {
		c.Cooldown = 5 * time.Second
	} else if c.RunFunc == nil {
		panic("RunFunc cannot be empty")
	}

	sort.Strings(c.Aliases)
}
