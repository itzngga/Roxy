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
	QuestionState        *skipmap.StringMap[*command.QuestionState]
	QuestionChan         chan *command.QuestionState
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

func (m *Muxer) HandleQuestionStateChan() {
	go func() {
		for message := range m.QuestionChan {
			for _, question := range message.Questions {
				if question.GetAnswer() == "" {
					message.ActiveQuestion = question.Question
					m.QuestionState.Store(message.RunFuncCtx.Number, message)
					message.RunFuncCtx.SendReplyMessage(question.Question)
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

func (m *Muxer) HandleQuestionState(c *whatsmeow.Client, evt *events.Message, parsedMsg string) bool {
	number := evt.Info.Sender.ToNonAD().String()
	questionState, ok := m.QuestionState.Load(number)
	if ok {
		if strings.Contains(parsedMsg, "cancel") || strings.Contains(parsedMsg, "batal") {
			m.QuestionState.Delete(number)
			return false
		} else {
			for i, question := range questionState.Questions {
				if question.Question == questionState.ActiveQuestion && question.GetAnswer() == "" {
					if questionState.Questions[i].Capture {
						questionState.Questions[i].SetAnswer(evt.Message)
					} else {
						questionState.Questions[i].SetAnswer(parsedMsg)
					}
					continue
				} else if question.Question != questionState.ActiveQuestion && question.GetAnswer() == "" {
					questionState.ActiveQuestion = question.Question
					go func() {
						jids := []waTypes.MessageID{
							evt.Info.ID,
						}
						c.MarkRead(jids, evt.Info.Timestamp, evt.Info.Chat, evt.Info.Sender)
					}()
					questionState.RunFuncCtx.SendReplyMessage(question.Question)
					return false
				} else if question.Question == questionState.ActiveQuestion && question.GetAnswer() != "" {
					continue
				}
			}

			m.QuestionState.Delete(number)
			questionState.ResultChan <- true
			return true
		}
	}
	return true
}

func (m *Muxer) RunCommand(c *whatsmeow.Client, evt *events.Message) {
	number := evt.Info.Sender.ToNonAD().String()
	if midOk := m.GlobalMiddlewareProcessing(c, evt, number); !midOk {
		return
	}
	parsed := util.ParseMessageText(evt)

	var fromMe = number == c.Store.ID.ToNonAD().String()
	if ok := m.HandleQuestionState(c, evt, parsed); !fromMe && ok {
		parsed = util.ParseMessageText(evt)
	} else if !fromMe && !ok {
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

		var args = strings.Split(parsed, " ")[1:]
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
			QuestionChan:   m.QuestionChan,
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
		QuestionState:        skipmap.NewString[*command.QuestionState](),
		Categories:           skipmap.NewString[string](),
		QuestionChan:         make(chan *command.QuestionState),
		MessageTimeout:       options.SendMessageTimeout,
		Options:              options,
		Log:                  log,
	}
	muxer.PrepareDefaultMiddleware()
	muxer.HandleQuestionStateChan()

	muxer.AddAllEmbed()

	return muxer
}
