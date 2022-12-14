package core

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/middleware"
	"github.com/itzngga/goRoxy/options"
	"github.com/itzngga/goRoxy/types"
	"github.com/itzngga/goRoxy/util"
	"github.com/zhangyunhao116/skipmap"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	types2 "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var Commands types.Embed[*command.Command]
var Middlewares types.Embed[command.MiddlewareFunc]
var Categories types.Embed[string]

type Muxer struct {
	Options        *options.Options
	Log            waLog.Logger
	MessageTimeout time.Duration

	Categories           *skipmap.StringMap[string]
	GlobalMiddleware     *skipmap.StringMap[command.MiddlewareFunc]
	Commands             *skipmap.StringMap[*command.Command]
	CommandResponseCache *skipmap.StringMap[*waProto.Message]
	Locals               *skipmap.StringMap[string]
}

func (m *Muxer) HandlePanic() {
	rec := recover()
	if rec != nil {
		m.Log.Errorf("panic: \n%v", string(debug.Stack()))
	}
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
	m.GlobalMiddleware.Range(func(key string, middleware command.MiddlewareFunc) bool {
		m.GlobalMiddleware.Delete(key)
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
	categories := Categories.Get()
	for _, cat := range categories {
		m.Categories.Store(cat, cat)
	}
	commands := Commands.Get()
	for _, cmd := range commands {
		m.AddCommand(cmd)
	}
	middlewares := Middlewares.Get()
	for _, mid := range middlewares {
		m.AddMiddleware(mid)
	}
}

func (m *Muxer) AddMiddleware(middleware command.MiddlewareFunc) {
	m.GlobalMiddleware.Store(uuid.New().String(), middleware)
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
	cmdLoad, isAvailable := m.Commands.Load(cmd)
	if isCmd && isAvailable {
		defer m.HandlePanic()
		args := command.RunFuncArgs{
			Options: m.Options,
			Evm:     evt,
			Cmd:     cmdLoad,
			Msg:     parsed,
			Number:  number,
			Locals:  m.Locals,
			Args:    strings.Split(parsed, " "),
		}
		midAreOk := true
		m.GlobalMiddleware.Range(func(key string, value command.MiddlewareFunc) bool {
			if !value(c, args) {
				midAreOk = false
				return false
			}
			return true
		})
		if !midAreOk {
			return
		}
		if cmdLoad.Middleware != nil {
			if !cmdLoad.Middleware(c, args) {
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
				msg = cmdLoad.RunFunc(c, args)
			}
		} else {
			msg = cmdLoad.RunFunc(c, args)
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
		jids := []types2.MessageID{
			evt.Info.ID,
		}
		c.MarkRead(jids, evt.Info.Timestamp, evt.Info.Chat, evt.Info.Sender)
	}
}

func (m *Muxer) PrepareDefaultMiddleware() {
	if m.Options.WithCommandLog {
		m.GlobalMiddleware.Store("log", middleware.LogMiddleware)
	}
	if m.Options.WithCommandCooldown {
		if m.Options.CommandCooldownTimeout == 0 {
			m.Options.CommandCooldownTimeout = 5
		}
		middleware.CooldownCache = skipmap.NewString[string]()
		m.GlobalMiddleware.Store("cooldown", middleware.CooldownMiddleware)
	}
}

func NewMuxer(log waLog.Logger, options *options.Options) *Muxer {
	muxer := &Muxer{
		Locals:               skipmap.NewString[string](),
		Commands:             skipmap.NewString[*command.Command](),
		GlobalMiddleware:     skipmap.NewString[command.MiddlewareFunc](),
		CommandResponseCache: skipmap.NewString[*waProto.Message](),
		Categories:           skipmap.NewString[string](),
		MessageTimeout:       options.SendMessageTimeout,
		Options:              options,
		Log:                  log,
	}
	muxer.PrepareDefaultMiddleware()
	muxer.Locals.Store("uid", util.CreateUid())

	muxer.AddAllEmbed()
	muxer.GenerateHelpButton()

	return muxer
}

func init() {
	cmd := types.NewEmbed[*command.Command]()
	mid := types.NewEmbed[command.MiddlewareFunc]()
	cat := types.NewEmbed[string]()
	Commands = &cmd
	Middlewares = &mid
	Categories = &cat
}
