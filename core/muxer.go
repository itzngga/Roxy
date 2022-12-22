package core

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	messageMiddlewares "github.com/itzngga/goRoxy/basic/message_middleware"
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/embed"
	"github.com/itzngga/goRoxy/options"
	"github.com/itzngga/goRoxy/util"
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
	m.GenerateHelpButton()
}

func (m *Muxer) AddAllEmbed() {
	categories := embed.Categories.Get()
	for _, cat := range categories {
		m.Categories.Store(cat, cat)
	}
	commands := embed.Commands.Get()
	for _, cmd := range commands {
		if !m.Options.WithBuiltIn {
			if cmd.BuiltIn {
				continue
			}
			m.AddCommand(cmd)
		} else {
			m.AddCommand(cmd)
		}

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

func (m *Muxer) CheckGlobalState(number string) (bool, string) {
	globalState, ok := m.Locals.Load(number)
	if !ok {
		return false, ""
	}
	m.Locals.Delete(number)
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
	param := &command.RunFuncParams{
		Options: m.Options,
		Event:   evt,
		Number:  number,
		Locals:  m.Locals,
	}
	midAreOk := true
	m.GlobalMiddlewares.Range(func(key string, value command.MiddlewareFunc) bool {
		if !value(c, param) {
			midAreOk = false
			return false
		}
		return true
	})

	return midAreOk
}

func (m *Muxer) RunCommand(c *whatsmeow.Client, evt *events.Message) {
	number := strconv.FormatUint(evt.Info.Sender.UserInt(), 10)
	if midOk := m.GlobalMiddlewareProcessing(c, evt, number); !midOk {
		return
	}
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
	cmdLoad, isAvailable := m.Commands.Load(cmd)
	if isCmd && isAvailable {
		params := &command.RunFuncParams{
			Log:       m.Log,
			Options:   m.Options,
			Event:     evt,
			Info:      &evt.Info,
			User:      c.Store.ID,
			Message:   evt.Message,
			Cmd:       cmdLoad,
			ParsedMsg: parsed,
			Number:    number,
			Locals:    m.Locals,
			Arguments: strings.Split(parsed, " "),
		}
		midAreOk := true
		m.Middlewares.Range(func(key string, value command.MiddlewareFunc) bool {
			if !value(c, params) {
				midAreOk = false
				return false
			}
			return true
		})
		if !midAreOk {
			return
		}
		if cmdLoad.Middleware != nil {
			if !cmdLoad.Middleware(c, params) {
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
			msg = m.GetCachedCommandResponse(cmdLoad.Name)
			if msg == nil {
				msg = cmdLoad.RunFunc(c, params)
			}
		} else {
			msg = cmdLoad.RunFunc(c, params)
		}
		if msg != nil {
			ctx, cancel := context.WithTimeout(context.Background(), m.MessageTimeout)
			defer cancel()

			_, err := c.SendMessage(ctx, evt.Info.Chat, "", msg)
			if err != nil {
				fmt.Println("[SEND MESSAGE ERR]\n", err)
			}
			if cmdLoad.Cache {
				m.SetCacheCommandResponse(cmdLoad.Name, msg)
			}
		}
		jids := []waTypes.MessageID{
			evt.Info.ID,
		}
		c.MarkRead(jids, evt.Info.Timestamp, evt.Info.Chat, evt.Info.Sender)
	}
}

func (m *Muxer) PrepareDefaultMiddleware() {
	if m.Options.WithCommandLog {
		m.Middlewares.Store("log", messageMiddlewares.LogMiddleware)
	}
	if m.Options.WithCommandCooldown {
		if m.Options.CommandCooldownTimeout == 0 {
			m.Options.CommandCooldownTimeout = 5
		}
		messageMiddlewares.CooldownCache = skipmap.NewString[string]()
		m.Middlewares.Store("cooldown", messageMiddlewares.CooldownMiddleware)
	}
}

func NewMuxer(log waLog.Logger, options *options.Options) *Muxer {
	muxer := &Muxer{
		Locals:               skipmap.NewString[string](),
		Commands:             skipmap.NewString[*command.Command](),
		GlobalMiddlewares:    skipmap.NewString[command.MiddlewareFunc](),
		Middlewares:          skipmap.NewString[command.MiddlewareFunc](),
		CommandResponseCache: skipmap.NewString[*waProto.Message](),
		Categories:           skipmap.NewString[string](),
		MessageTimeout:       options.SendMessageTimeout,
		Options:              options,
		Log:                  log,
	}
	muxer.PrepareDefaultMiddleware()
	muxer.Locals.Store("uid", util.CreateUid())

	muxer.AddAllEmbed()
	if options.WithHelpCommand {
		muxer.GenerateHelpButton()
	}

	return muxer
}
