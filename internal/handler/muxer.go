package handler

import (
	"fmt"
	"github.com/itzngga/goRoxy/helper"
	"github.com/itzngga/goRoxy/util"
	"github.com/zhangyunhao116/skipmap"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"strconv"
	"strings"
)

var GlobalMiddleware *skipmap.StringMap[MiddlewareFunc]

type Muxer struct {
	CmdMap *skipmap.StringMap[*Command]
	Locals *skipmap.StringMap[string]
}

func (m *Muxer) AddCommand(cmd *Command) {
	cmd.Validate()
	_, ok := m.CmdMap.Load(cmd.Name)
	if ok {
		panic("Duplicate command: " + cmd.Name)
	}

	for _, alias := range cmd.Aliases {
		_, ok := m.CmdMap.Load(alias)
		if ok {
			panic("Duplicate alias in command " + cmd.Name)
		}
		m.CmdMap.Store(alias, cmd)
	}
	m.CmdMap.Store(cmd.Name, cmd)
}

func (m *Muxer) CheckGlobalState(number string) (bool, string) {
	globalState, ok := m.Locals.Load(number)
	if !ok {
		return false, ""
	}
	m.Locals.Delete(number)
	return true, globalState
}

func (m *Muxer) RunCommand(c *whatsmeow.Client, evt *events.Message) {
	number := strconv.FormatUint(evt.Info.Sender.UserInt(), 10)
	id, _ := m.Locals.Load("uid")
	ok, stateCmd := m.CheckGlobalState(number)
	parsed := util.ParseMessageText(id, evt)
	if ok {
		if strings.Contains(parsed, "!cancel") {
			m.Locals.Delete(number)
			util.SendReplyMessage(c, evt, "Dibatalkan!")
			return
		}
		parsed = stateCmd + " " + parsed
	}
	cmd, isCmd := util.ParseCmd(parsed)
	command, isAvailable := m.CmdMap.Load(cmd)
	if isCmd && isAvailable {
		args := RunFuncArgs{
			Evm:    evt,
			Cmd:    command,
			Msg:    parsed,
			Number: number,
			Locals: m.Locals,
			Args:   strings.Split(parsed, " "),
		}
		GlobalMiddleware.Range(func(key string, value MiddlewareFunc) bool {
			if !value(c, args) {
				return false
			}
			return false
		})

		if command.Middleware != nil {
			if !command.Middleware(c, args) {
				return
			}
		}
		if command.GroupOnly {
			if !evt.Info.IsGroup {
				return
			}
		}
		if command.PrivateOnly {
			if evt.Info.IsGroup {
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
		a := []types.MessageID{}
		a = append(a, evt.Info.ID)
		c.MarkRead(a, evt.Info.Timestamp, evt.Info.Chat, evt.Info.Sender)
	}
}

func (m *Muxer) GenerateRequiredLocals() {
	uid := helper.CreateUid()
	m.Locals.Store("uid", uid)
}

func NewMuxer() *Muxer {
	muxer := &Muxer{
		Locals: skipmap.NewString[string](),
		CmdMap: skipmap.NewString[*Command](),
	}
	muxer.GenerateRequiredLocals()
	return muxer
}
