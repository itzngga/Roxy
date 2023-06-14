package core

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/options"
	"github.com/itzngga/Roxy/util"
	"github.com/zhangyunhao116/skipmap"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"strconv"
	"strings"
	"time"
)

type Muxer struct {
	Options              *options.Options
	Log                  waLog.Logger
	MessageTimeout       time.Duration
	Categories           *skipmap.StringMap[string]
	GlobalMiddlewares    *skipmap.StringMap[command.MiddlewareFunc]
	Middlewares          *skipmap.StringMap[command.MiddlewareFunc]
	Commands             *skipmap.StringMap[*command.Command]
	CommandResponseCache *skipmap.StringMap[*waProto.Message]
	UserState            *skipmap.StringMap[*command.StateCommand]
	UserStateChan        chan []interface{}
	Locals               *skipmap.StringMap[string]
}

func (m *Muxer) Clean() {
	m.Categories.Range(func(key string, category string) bool {
		m.Categories.Delete(key)
		return true
	})
	m.Commands.Range(func(key string, cmd *command.Command) bool {
		m.Commands.Delete(key)
		return true
	})
	m.Middlewares.Range(func(key string, middleware command.MiddlewareFunc) bool {
		m.Middlewares.Delete(key)
		return true
	})
	m.Locals.Range(func(key string, middleware string) bool {
		m.Locals.Delete(key)
		return true
	})
	m.PrepareDefaultMiddleware()
	m.Locals.Store("uid", util.CreateUid())
}

func (m *Muxer) HandleUserStateChannel() {
	go func() {
		for message := range m.UserStateChan {
			for _, state := range embed.StateCommand.Get() {
				if state.Name == message[0] {
					state.Locals = message[2].(map[string]interface{})
					m.UserState.Store(message[1].(string), state)
					go func() {
						timeout := time.NewTimer(state.StateTimeout)
						<-timeout.C
						m.UserState.Delete(message[1].(string))
						timeout.Stop()
					}()
					break
				}
			}
		}
	}()
}

func (m *Muxer) AddAllEmbed() {
	categories := embed.Categories.Get()
	for _, cat := range categories {
		m.Categories.Store(cat, cat)
	}
	commands := embed.Commands.Get()
	for _, cmd := range commands {
		m.AddCommand(cmd)
	}
	middlewares := embed.Middlewares.Get()
	for _, mid := range middlewares {
		m.AddMiddleware(mid)

	}
	globalMiddleware := embed.GlobalMiddlewares.Get()
	for _, mid := range globalMiddleware {
		m.AddGlobalMiddleware(mid)
	}
}

func (m *Muxer) AddGlobalMiddleware(middleware command.MiddlewareFunc) {
	m.GlobalMiddlewares.Store(uuid.New().String(), middleware)
}

func (m *Muxer) AddMiddleware(middleware command.MiddlewareFunc) {
	m.Middlewares.Store(uuid.New().String(), middleware)
}

func (m *Muxer) AddCommand(cmd *command.Command) {
	cmd.Validate()
	_, ok := m.Commands.Load(cmd.Name)
	if ok {
		panic("error: duplicate command " + cmd.Name)
	}

	for _, alias := range cmd.Aliases {
		_, ok := m.Commands.Load(alias)
		if ok {
			panic("error: duplicate alias in command " + cmd.Name)
		}
		m.Commands.Store(alias, cmd)
	}
	m.Commands.Store(cmd.Name, cmd)
}

func (m *Muxer) CheckGlobalState(number string) (bool, *command.StateCommand) {
	globalState, ok := m.UserState.Load(number)
	if !ok {
		return false, nil
	}
	return true, globalState
}

func (m *Muxer) GetCachedCommandResponse(cmd string) *waProto.Message {
	cache, ok := m.CommandResponseCache.Load(cmd)
	if ok {
		return cache
	}
	return nil
}

func (m *Muxer) SetCacheCommandResponse(cmd string, response *waProto.Message) {
	m.CommandResponseCache.Store(cmd, response)
	go func() {
		timeout := time.NewTimer(m.Options.CommandResponseCacheTimeout)
		<-timeout.C
		m.CommandResponseCache.Delete(cmd)
		timeout.Stop()
	}()
}

func (m *Muxer) GlobalMiddlewareProcessing(c *whatsmeow.Client, evt *events.Message, number string) bool {
	param := &command.RunFuncContext{
		Client:       c,
		Options:      m.Options,
		MessageEvent: evt,
		Number:       number,
		Locals:       m.Locals,
	}
	midAreOk := true
	m.GlobalMiddlewares.Range(func(key string, value command.MiddlewareFunc) bool {
		if !value(param) {
			midAreOk = false
			return false
		}
		return true
	})

	return midAreOk
}

func (m *Muxer) HandleUserState(c *whatsmeow.Client, evt *events.Message, parsedMsg string) bool {
	number := evt.Info.Sender.ToNonAD().String()
	ok, stateCmd := m.CheckGlobalState(number)
	if ok {
		if stateCmd.PrivateOnly {
			if evt.Info.IsGroup {
				return true
			}
		}

		if stateCmd.GroupOnly {
			if !evt.Info.IsGroup {
				return true
			}
		}

		params := &command.StateFuncContext{
			Client:       c,
			WaLog:        m.Log,
			Options:      m.Options,
			MessageEvent: evt,
			MessageInfo:  &evt.Info,
			ClientJID:    c.Store.ID,
			Message:      evt.Message,
			CurrentState: stateCmd,
			ParsedMsg:    parsedMsg,
			Number:       number,
			Locals:       stateCmd.Locals,
		}

		if strings.Contains(parsedMsg, "cancel") {
			params.SendReplyMessage(stateCmd.CancelReply)
		} else if strings.Contains(parsedMsg, "batal") {
			params.SendReplyMessage(stateCmd.CancelReply)
		}

		stateCmd.RunFunc(params)
		return false
	}
	return true
}

func (m *Muxer) RunCommand(c *whatsmeow.Client, evt *events.Message) {
	number := strconv.FormatUint(evt.Info.Sender.UserInt(), 10)
	if midOk := m.GlobalMiddlewareProcessing(c, evt, number); !midOk {
		return
	}
	id, _ := m.Locals.Load("uid")
	parsed := util.ParseMessageText(id, evt)

	if ok := m.HandleUserState(c, evt, parsed); !ok {
		return
	}

	prefix, cmd, isCmd := util.ParseCmd(parsed)
	cmdLoad, isAvailable := m.Commands.Load(cmd)
	if isCmd && isAvailable {
		go func() {
			jids := []waTypes.MessageID{
				evt.Info.ID,
			}
			c.MarkRead(jids, evt.Info.Timestamp, evt.Info.Chat, evt.Info.Sender)
		}()
		var fromMe bool
		if id = evt.Info.Sender.ToNonAD().String(); *util.ParseQuotedRemoteJid(evt) == id {
			fromMe = true
		}
		var args = strings.Split(parsed, " ")
		params := &command.RunFuncContext{
			Client:         c,
			WaLog:          m.Log,
			Options:        m.Options,
			MessageEvent:   evt,
			MessageInfo:    &evt.Info,
			ClientJID:      c.Store.ID,
			Message:        evt.Message,
			FromMe:         fromMe,
			CurrentCommand: cmdLoad,
			ParsedMsg:      parsed,
			Number:         number,
			Locals:         m.Locals,
			UserStateChan:  m.UserStateChan,
			Prefix:         prefix,
			Arguments:      args,
		}
		var midAreOk = true
		m.Middlewares.Range(func(key string, value command.MiddlewareFunc) bool {
			if !value(params) {
				midAreOk = false
				return false
			}
			return true
		})
		if !midAreOk {
			return
		}
		if cmdLoad.Middleware != nil {
			if !cmdLoad.Middleware(params) {
				return
			}
		}
		if cmdLoad.GroupOnly {
			if !evt.Info.IsGroup {
				return
			}
		}
		if cmdLoad.PrivateOnly {
			if evt.Info.IsGroup {
				return
			}
		}
		var msg *waProto.Message
		if cmdLoad.Cache {
			msg = m.GetCachedCommandResponse(parsed)
			if msg == nil {
				msg = cmdLoad.RunFunc(params)
			}
		} else {
			msg = cmdLoad.RunFunc(params)
		}
		if msg != nil {
			ctx, cancel := context.WithTimeout(context.Background(), m.Options.SendMessageTimeout)
			defer cancel()

			_, err := c.SendMessage(ctx, evt.Info.Chat, msg)
			if err != nil {
				fmt.Println("[SEND MESSAGE ERR]\n", err)
			}
			if cmdLoad.Cache {
				m.SetCacheCommandResponse(parsed, msg)
			}
		}
	}
}

func (m *Muxer) PrepareDefaultMiddleware() {
	if m.Options.WithCommandLog {
		m.Middlewares.Store("log", func(ctx *command.RunFuncContext) bool {
			ctx.WaLog.Infof("[CMD] [%s] command > %s", ctx.Number, ctx.CurrentCommand.Name)
			return true
		})
	}
}

func NewMuxer(log waLog.Logger, options *options.Options) *Muxer {
	muxer := &Muxer{
		Locals:               skipmap.NewString[string](),
		Commands:             skipmap.NewString[*command.Command](),
		GlobalMiddlewares:    skipmap.NewString[command.MiddlewareFunc](),
		Middlewares:          skipmap.NewString[command.MiddlewareFunc](),
		CommandResponseCache: skipmap.NewString[*waProto.Message](),
		UserState:            skipmap.NewString[*command.StateCommand](),
		Categories:           skipmap.NewString[string](),
		UserStateChan:        make(chan []interface{}),
		MessageTimeout:       options.SendMessageTimeout,
		Options:              options,
		Log:                  log,
	}
	muxer.PrepareDefaultMiddleware()
	muxer.HandleUserStateChannel()
	muxer.Locals.Store("uid", util.CreateUid())

	muxer.AddAllEmbed()

	return muxer
}
