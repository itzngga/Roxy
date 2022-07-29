package handler

import (
	"fmt"
	"sort"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type MiddlewareFunc func(c *whatsmeow.Client, m *events.Message) bool
type RunFunc func(c *whatsmeow.Client, m *events.Message) *waProto.Message

type Category struct {
	Name        string
	Description string
}

var (
	UtilitiesCategory = &Category{
		Name:        "Utilities",
		Description: "Bot Utilities",
	}
	MiscCategory = &Category{
		Name:        "Misc",
		Description: "Bot Misc",
	}
	Uncategorized = &Category{
		Name:        "Uncategorized",
		Description: "Bot Uncategorized",
	}
)

type Command struct {
	Name            string
	Aliases         []string
	Description     string
	LongDescription string

	Cooldown time.Duration
	Category *Category

	HideFromHelp bool
	GroupOnly    bool
	Middleware   MiddlewareFunc
	RunFunc      RunFunc
}

func SendInvalidCommand(m *events.Message) *waProto.Message {
	return &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String("Invalid command received"),
			ContextInfo: WithReply(m),
		},
	}
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

func AddCommand(cmd *Command) {
	DefaultMuxer.AddCommand(cmd)
}

func RunCommand(c *whatsmeow.Client, evt *events.Message) {
	go DefaultMuxer.RunCommand(c, evt)
}
