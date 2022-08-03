package handler

import (
	"fmt"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"sync"
	"time"

	"github.com/itzngga/goRoxy/internal/dictpool"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

var DefaultMuxer *Muxer

type Muxer struct {
	mutex          *sync.RWMutex
	CommandSlice   []*Command
	DictPool       *dictpool.Dict
	HelpString     string
	CommandSucceed int
}

func (m *Muxer) FindCommand(cmd string) (*Command, bool) {
	i := m.DictPool.Get(cmd)
	if i != nil && m.CommandSlice[i.(int)] != nil {
		return m.CommandSlice[i.(int)], true
	} else {
		m.mutex.RLock()
		defer m.mutex.RUnlock()

		for i, val := range m.CommandSlice {
			if pcmd := val.GetName(cmd); pcmd == cmd {
				m.DictPool.Set(pcmd, i)
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
	i := m.DictPool.Get(cmd.Name)
	if i != nil && m.CommandSlice[i.(int)] != nil {
		panic(fmt.Sprintf("Invalid duplicate %s command", cmd.Name))
	}

	m.DictPool.Set(cmd.Name, len(m.CommandSlice))
	m.CommandSlice = append(m.CommandSlice, cmd)
}

func (m *Muxer) RunCommand(c *whatsmeow.Client, evt *events.Message) {
	parsed := ParseMessageText(evt)
	cmd, isCmd := ParseCmd(parsed)
	command, isAvailable := m.FindCommand(cmd)
	if isCmd && isAvailable {
		cdId := m.DictPool.Get(c.Store.ID.User + evt.Info.Sender.User)
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

		m.mutex.RLock()

		command.CommandSucceed = command.CommandSucceed + 1
		m.DictPool.Set(c.Store.ID.User+evt.Info.Sender.User, true)
		defer m.mutex.RUnlock()

		time.Sleep(5 * time.Second)
		m.DictPool.Del(c.Store.ID.User + evt.Info.Sender.User)

		return
	}
}

func (m *Muxer) UpdateHelp() {
	if m.HelpString == "" {
		m.HelpString = fmt.Sprintf(`---------------%s`, "\nTestBot\n\n")
	}
	var (
		Uncategorizeds []*Command
		Utilities      []*Command
		Misc           []*Command
	)
	for _, val := range m.CommandSlice {
		if !val.HideFromHelp {
			if val.Category == UtilitiesCategory {
				Utilities = append(Utilities, val)
			} else if val.Category == MiscCategory {
				Misc = append(Misc, val)
			} else {
				Uncategorizeds = append(Uncategorizeds, val)
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
	for _, unca := range Uncategorizeds {
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

func NewMuxer() *Muxer {
	muxer := &Muxer{
		mutex:        &sync.RWMutex{},
		DictPool:     dictpool.AcquireDict(),
		CommandSlice: make([]*Command, 0),
	}
	muxer.AddCommand(&Command{
		Name:        "help",
		Description: "Returns Bot Help",
		Category:    UtilitiesCategory,
		RunFunc: func(c *whatsmeow.Client, m *events.Message) *waProto.Message {
			return SendReplyText(m, muxer.GetHelpPage())
		},
	})
	return muxer
}

func NewDefaultMuxer() {
	DefaultMuxer = NewMuxer()
}
