package handler

import (
	"fmt"
	"github.com/itzngga/goRoxy/helper"
	"github.com/itzngga/goRoxy/util"
	"github.com/jellydator/ttlcache/v2"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"strings"

	"sync"
)

var GlobalLocals *map[string]interface{}
var GlobalMiddleware []MiddlewareFunc

type Muxer struct {
	mutex        *sync.RWMutex
	CommandSlice []*Command
	MuxCache     ttlcache.SimpleCache
	Locals       *map[string]interface{}
	HelpString   string
}

func (m *Muxer) FindCommand(cmd string) (*Command, bool) {
	i, _ := m.MuxCache.Get(cmd)
	if i != nil && len(m.CommandSlice) != 0 && m.CommandSlice[i.(int)] != nil {
		return m.CommandSlice[i.(int)], true
	} else {
		for i, val := range m.CommandSlice {
			if pcmd := val.GetName(cmd); pcmd == cmd {
				m.MuxCache.Set(pcmd, i)
				return val, true
			}
		}
		return &Command{}, false
	}
}

func (m *Muxer) AddCommand(cmd *Command) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	cmd.Validate()
	cmd.Locals = m.Locals

	indexCache, _ := m.MuxCache.Get(cmd.Name)
	if indexCache != nil && m.CommandSlice[indexCache.(int)] != nil {
		panic(fmt.Sprintf("Invalid duplicate %s cmd", cmd.Name))
	}

	m.MuxCache.Set(cmd.Name, len(m.CommandSlice))
	m.CommandSlice = append(m.CommandSlice, cmd)
}

func (m *Muxer) GetLocals(key string) interface{} {
	return (*m.Locals)[key]
}

func (m *Muxer) SetLocals(key string, value interface{}) {
	(*m.Locals)[key] = value
}

func (m *Muxer) RunCommand(c *whatsmeow.Client, evt *events.Message) {
	parsed := util.ParseMessageText(m.GetLocals("uid").(string), evt)
	cmd, isCmd := util.ParseCmd(parsed)
	command, isAvailable := m.FindCommand(cmd)
	if isCmd && isAvailable {
		args := RunFuncArgs{
			Evm:  evt,
			Cmd:  command,
			Msg:  parsed,
			Args: strings.Split(parsed, " "),
		}
		for _, middlewareFunc := range GlobalMiddleware {
			if middlewareFunc != nil {
				if m := middlewareFunc(c, args); !m {
					return
				}
			}
		}
		if command.Middleware != nil {
			if m := command.Middleware(c, args); !m {
				return
			}
		}
		if command.GroupOnly {
			if !evt.Info.IsGroup {
				return
			}
		}
		msg := command.RunFunc(c, args)

		if msg != nil {
			_, err := c.SendMessage(evt.Info.Chat, "", msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (m *Muxer) UpdateHelp() {
	if m.HelpString == "" {
		m.HelpString = fmt.Sprintf(`---------------%s`, "\nTestBot\n\n")
	}
	var Uncategorize []*Command
	var Utilities []*Command
	var Misc []*Command
	for _, val := range m.CommandSlice {
		if !val.HideFromHelp {
			if val.Category == UtilitiesCategory {
				Utilities = append(Utilities, val)
			} else if val.Category == MiscCategory {
				Misc = append(Misc, val)
			} else {
				Uncategorize = append(Uncategorize, val)
			}
		}
	}
	m.HelpString += "#Utilities\n" + UtilitiesCategory.Description + "\n"
	for _, util := range Utilities {
		m.HelpString += fmt.Sprintf(`%s. %s%s`, "➣ ", util.Name+"\n", util.Description+"\n")
	}
	m.HelpString += "\n\n#Misc\n" + MiscCategory.Description + "\n"
	for _, misc := range Misc {
		m.HelpString += fmt.Sprintf(`%s. %s%s`, "➣ ", misc.Name+"\n", misc.Description+"\n")
	}
	m.HelpString += "\n\n#Uncategorized\n" + Uncategorized.Description + "\n"
	for _, unca := range Uncategorize {
		m.HelpString += fmt.Sprintf(`%s. %s%s`, "➣ ", unca.Name+"\n", unca.Description+"\n")
	}
}

func (m *Muxer) GetHelpPage() string {
	if m.HelpString != "" {
		return m.HelpString
	}
	m.UpdateHelp()
	return m.HelpString
}

func (m *Muxer) GenerateRequiredLocals() {
	uid := helper.CreateUid()
	m.SetLocals("uid", uid)
}

func NewMuxer() *Muxer {
	muxer := &Muxer{
		mutex:        &sync.RWMutex{},
		MuxCache:     ttlcache.NewCache(),
		Locals:       &map[string]interface{}{},
		CommandSlice: make([]*Command, 0),
	}
	muxer.GenerateRequiredLocals()
	muxer.AddCommand(&Command{
		Name:        "help",
		Description: "Returns Bot Help",
		Category:    UtilitiesCategory,
		RunFunc: func(c *whatsmeow.Client, args RunFuncArgs) *waProto.Message {
			return util.SendReplyText(args.Evm, muxer.GetHelpPage())
		},
	})
	return muxer
}
