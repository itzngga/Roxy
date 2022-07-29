package handler

import (
	"fmt"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v2"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

var DefaultMuxer = &defaultMuxer
var defaultMuxer Muxer

type Muxer struct {
	mutex          *sync.RWMutex
	CommandSlice   []*Command
	MuxCache       ttlcache.SimpleCache
	CooldownIndex  []int
	HelpString     string
	CommandSucceed int
}

func (m *Muxer) FindCommand(cmd string) (*Command, bool) {
	i, _ := m.MuxCache.Get(cmd)
	if i != nil && m.CommandSlice[i.(int)] != nil {
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
	indexCache, _ := m.MuxCache.Get(cmd.Name)
	if indexCache != nil && m.CommandSlice[indexCache.(int)] != nil {
		panic(fmt.Sprintf("Invalid duplicate %s command", cmd.Name))
	}

	m.MuxCache.Set(cmd.Name, len(m.CommandSlice))
	m.CommandSlice = append(m.CommandSlice, cmd)
}

func (m *Muxer) RunCommand(c *whatsmeow.Client, evt *events.Message) {
	parsed := ParseMessageText(evt)
	cmd, isCmd := ParseCmd(parsed)
	command, isAvailable := m.FindCommand(cmd)
	if isCmd && isAvailable {
		cdId, _ := m.MuxCache.Get(c.Store.ID.User + evt.Info.Sender.User)
		if cdId != nil {
			SendReplyMessage(c, evt, "You are on Cooldown!")
			return
		}
		if command.Middleware != nil {
			if m := command.Middleware(c, evt); !m {
				return
			}
		}
		if command.GroupOnly {
			if !evt.Info.IsGroup {
				return
			}
		}
		msg := command.RunFunc(c, evt)
		if msg != nil {
			c.SendMessage(evt.Info.Chat, "", msg)
		}

		go func() {
			m.mutex.RLock()
			defer m.mutex.RUnlock()

			m.CommandSucceed++
			m.MuxCache.Set(c.Store.ID.User+evt.Info.Sender.User, true)
			time.Sleep(5 * time.Second)
			m.MuxCache.Remove(c.Store.ID.User + evt.Info.Sender.User)
		}()
		return
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

func (m Muxer) GetHelpPage() string {
	if m.HelpString != "" {
		return m.HelpString
	}
	m.UpdateHelp()
	return m.HelpString
}

func NewDefaultMuxer() {
	defaultMuxer = Muxer{
		mutex:         &sync.RWMutex{},
		MuxCache:      ttlcache.NewCache(),
		CooldownIndex: []int{},
		CommandSlice:  []*Command{},
	}
	AddCommand(&Command{
		Name:        "help",
		Description: "Returns Bot Help",
		Category:    UtilitiesCategory,
		RunFunc: func(c *whatsmeow.Client, m *events.Message) *waProto.Message {
			return SendReplyText(m, defaultMuxer.GetHelpPage())
		},
	})
}
