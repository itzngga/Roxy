package core

import (
	"context"
	"fmt"
	"github.com/alitto/pond"
	"github.com/google/uuid"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/options"
	"github.com/itzngga/Roxy/types"
	"github.com/itzngga/Roxy/util"
	"github.com/sajari/fuzzy"
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
	GroupCache           *skipmap.StringMap[[]*waTypes.GroupInfo]

	QuestionChan    chan *command.QuestionState
	Locals          *skipmap.StringMap[string]
	SuggestionModel *fuzzy.Model
	ctx             *skipmap.StringMap[types.RoxyContext]
}

func NewMuxer(ctx *skipmap.StringMap[types.RoxyContext], log waLog.Logger, options *options.Options) *Muxer {
	muxer := &Muxer{
		Locals:               skipmap.NewString[string](),
		Commands:             skipmap.NewString[*command.Command](),
		GlobalMiddlewares:    skipmap.NewString[command.MiddlewareFunc](),
		Middlewares:          skipmap.NewString[command.MiddlewareFunc](),
		CommandResponseCache: skipmap.NewString[*waProto.Message](),
		QuestionState:        skipmap.NewString[*command.QuestionState](),
		Categories:           skipmap.NewString[string](),
		GroupCache:           skipmap.NewString[[]*waTypes.GroupInfo](),
		QuestionChan:         make(chan *command.QuestionState),
		MessageTimeout:       options.SendMessageTimeout,
		Options:              options,
		Log:                  log,
	}
	muxer.extendContext(ctx)
	muxer.handleQuestionStateChan()
	muxer.addAllEmbed()

	if options.CommandSuggestion {
		muxer.generateSuggestionModel()
	}

	return muxer
}

func (muxer *Muxer) Clean() {
	muxer.Categories.Range(func(key string, category string) bool {
		muxer.Categories.Delete(key)
		return true
	})
	muxer.Commands.Range(func(key string, cmd *command.Command) bool {
		muxer.Commands.Delete(key)
		return true
	})
	muxer.Middlewares.Range(func(key string, middleware command.MiddlewareFunc) bool {
		muxer.Middlewares.Delete(key)
		return true
	})
	muxer.Locals.Range(func(key string, middleware string) bool {
		muxer.Locals.Delete(key)
		return true
	})
}

func (muxer *Muxer) handleQuestionStateChan() {
	muxer.getPool().Submit(func() {
		for message := range muxer.QuestionChan {
			muxer.QuestionState.Delete(message.RunFuncCtx.Number)
			for _, question := range message.Questions {
				if question.GetAnswer() == "" {
					message.ActiveQuestion = question.Question
					muxer.QuestionState.Store(message.RunFuncCtx.Number, message)
					if question.Question != "" {
						message.RunFuncCtx.SendReplyMessage(question.Question)
					}
					break
				}
			}
		}
	})
}

func (muxer *Muxer) addAllEmbed() {
	categories := embed.Categories.Get()
	for _, cat := range categories {
		muxer.Categories.Store(cat, cat)
	}
	commands := embed.Commands.Get()
	for _, cmd := range commands {
		muxer.AddCommand(cmd)
	}
	middlewares := embed.Middlewares.Get()
	for _, mid := range middlewares {
		muxer.AddMiddleware(mid)

	}
	globalMiddleware := embed.GlobalMiddlewares.Get()
	for _, mid := range globalMiddleware {
		muxer.AddGlobalMiddleware(mid)
	}
}

func (muxer *Muxer) AddGlobalMiddleware(middleware command.MiddlewareFunc) {
	muxer.GlobalMiddlewares.Store(uuid.New().String(), middleware)
}

func (muxer *Muxer) AddMiddleware(middleware command.MiddlewareFunc) {
	muxer.Middlewares.Store(uuid.New().String(), middleware)
}

func (muxer *Muxer) AddCommand(cmd *command.Command) {
	cmd.Validate()
	_, ok := muxer.Commands.Load(cmd.Name)
	if ok {
		panic("error: duplicate command " + cmd.Name)
	}

	for _, alias := range cmd.Aliases {
		_, ok := muxer.Commands.Load(alias)
		if ok {
			panic("error: duplicate alias in command " + cmd.Name)
		}
		muxer.Commands.Store(alias, cmd)
	}
	muxer.Commands.Store(cmd.Name, cmd)
}

func (muxer *Muxer) getCachedCommandResponse(cmd string) *waProto.Message {
	cache, ok := muxer.CommandResponseCache.Load(cmd)
	if ok {
		return cache
	}
	return nil
}

func (muxer *Muxer) setCacheCommandResponse(cmd string, response *waProto.Message) {
	muxer.CommandResponseCache.Store(cmd, response)
	muxer.getPool().Submit(func() {
		timeout := time.NewTimer(muxer.Options.CommandResponseCacheTimeout)
		<-timeout.C
		muxer.CommandResponseCache.Delete(cmd)
		timeout.Stop()
	})
}

func (muxer *Muxer) globalMiddlewareProcessing(c *whatsmeow.Client, evt *events.Message, number string) bool {
	if muxer.GlobalMiddlewares.Len() >= 1 {
		param := &command.RunFuncContext{
			Client:       c,
			WaLog:        muxer.Log,
			Options:      muxer.Options,
			MessageEvent: evt,
			MessageInfo:  &evt.Info,
			ClientJID:    c.Store.ID,
			Message:      evt.Message,
			FromMe:       evt.Info.IsFromMe,
			Number:       number,
			Locals:       muxer.Locals,
			QuestionChan: muxer.QuestionChan,
			Ctx:          muxer.ctx,
		}

		midAreOk := true
		muxer.GlobalMiddlewares.Range(func(key string, value command.MiddlewareFunc) bool {
			if !value(param) {
				midAreOk = false
				return false
			}
			return true
		})
		return midAreOk
	}

	return true
}

func (muxer *Muxer) handleQuestionState(c *whatsmeow.Client, evt *events.Message, number, parsedMsg string) {
	questionState, _ := muxer.QuestionState.Load(number)
	if strings.Contains(parsedMsg, "cancel") || strings.Contains(parsedMsg, "batal") {
		muxer.QuestionState.Delete(number)
		return
	} else {
		muxer.getPool().Submit(func() {
			jids := []waTypes.MessageID{
				evt.Info.ID,
			}
			c.MarkRead(jids, evt.Info.Timestamp, evt.Info.Chat, evt.Info.Sender)
		})
		for i, question := range questionState.Questions {
			if question.Question == questionState.ActiveQuestion && question.GetAnswer() == "" {
				if questionState.Questions[i].Capture {
					questionState.Questions[i].SetAnswer(evt.Message)
				} else if questionState.Questions[i].Reply {
					result := util.GetQuotedText(evt)
					questionState.Questions[i].SetAnswer(result)
				} else {
					questionState.Questions[i].SetAnswer(parsedMsg)
				}
				if questionState.WithEmojiReact {
					util.SendEmojiMessage(c, evt, questionState.EmojiReact)
				}
				continue
			} else if question.Question != questionState.ActiveQuestion && question.GetAnswer() == "" {
				questionState.ActiveQuestion = question.Question
				if question.Question != "" {
					questionState.RunFuncCtx.SendReplyMessage(question.Question)
				}
				return
			} else if question.Question == questionState.ActiveQuestion && question.GetAnswer() != "" {
				continue
			}
		}

		muxer.QuestionState.Delete(number)
		questionState.ResultChan <- true
		return
	}
}

func (muxer *Muxer) RunCommand(c *whatsmeow.Client, evt *events.Message) {
	if muxer.Options.AllowFromPrivate && !muxer.Options.AllowFromGroup && evt.Info.IsGroup {
		return
	}
	if muxer.Options.AllowFromGroup && !muxer.Options.AllowFromPrivate && !evt.Info.IsGroup {
		return
	}
	if !muxer.Options.AllowFromGroup && !muxer.Options.AllowFromPrivate {
		return
	}
	if muxer.Options.OnlyFromSelf && !evt.Info.IsFromMe {
		return
	}

	number := evt.Info.Sender.ToNonAD().String()
	parsed := util.ParseMessageText(evt)
	if !evt.Info.IsFromMe {
		_, ok := muxer.QuestionState.Load(number)
		if ok {
			muxer.handleQuestionState(c, evt, number, parsed)
			return
		}
	}

	if midOk := muxer.globalMiddlewareProcessing(c, evt, number); !midOk {
		return
	}

	prefix, cmd, isCmd := util.ParseCmd(parsed)
	cmdLoad, isAvailable := muxer.Commands.Load(cmd)
	if muxer.Options.CommandSuggestion && isCmd && !isAvailable {
		muxer.suggestCommand(evt, prefix, cmd)
		return
	}

	if isCmd && isAvailable {
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
		if cmdLoad.OnlyAdminGroup && evt.Info.IsGroup {
			if ok, _ := muxer.isGroupAdmin(evt.Info.Chat, evt.Info.Sender); !ok {
				return
			}
		}
		if cmdLoad.OnlyIfBotAdmin && evt.Info.IsGroup {
			if ok, _ := muxer.isClientGroupAdmin(evt.Info.Chat); !ok {
				return
			}
		}
		muxer.getPool().Submit(func() {
			if muxer.Options.WithCommandLog {
				muxer.Log.Infof("[CMD] [%s] command > %s", number, cmdLoad.Name)
			}
			jids := []waTypes.MessageID{
				evt.Info.ID,
			}
			c.MarkRead(jids, evt.Info.Timestamp, evt.Info.Chat, evt.Info.Sender)
		})

		var args = strings.Split(parsed, " ")[1:]
		params := &command.RunFuncContext{
			Client:         c,
			WaLog:          muxer.Log,
			Options:        muxer.Options,
			MessageEvent:   evt,
			MessageInfo:    &evt.Info,
			ClientJID:      c.Store.ID,
			Message:        evt.Message,
			FromMe:         evt.Info.IsFromMe,
			CurrentCommand: cmdLoad,
			ParsedMsg:      parsed,
			Number:         number,
			Locals:         muxer.Locals,
			QuestionChan:   muxer.QuestionChan,
			Prefix:         prefix,
			Arguments:      args,
			Ctx:            muxer.ctx,
		}
		if muxer.Middlewares.Len() >= 1 {
			var midAreOk = true
			muxer.Middlewares.Range(func(key string, value command.MiddlewareFunc) bool {
				if !value(params) {
					midAreOk = false
					return false
				}
				return true
			})
			if !midAreOk {
				return
			}
		}
		if cmdLoad.Middleware != nil {
			if !cmdLoad.Middleware(params) {
				return
			}
		}
		var msg *waProto.Message
		if cmdLoad.Cache {
			msg = muxer.getCachedCommandResponse(parsed)
			if msg == nil {
				msg = cmdLoad.RunFunc(params)
			}
		} else {
			msg = cmdLoad.RunFunc(params)
		}
		if msg != nil {
			ctx, cancel := context.WithTimeout(context.Background(), muxer.Options.SendMessageTimeout)
			defer cancel()

			_, err := c.SendMessage(ctx, evt.Info.Chat, msg)
			if err != nil {
				fmt.Println("[SEND MESSAGE ERR]\n", err)
			}
			if cmdLoad.Cache {
				muxer.setCacheCommandResponse(parsed, msg)
			}
		}
	}
}

func (muxer *Muxer) getPool() *pond.WorkerPool {
	return types.GetContext[*pond.WorkerPool](muxer.ctx, "workerPool")
}

func (muxer *Muxer) getCurrentClient() *whatsmeow.Client {
	return types.GetContext[*whatsmeow.Client](muxer.ctx, "appClient")
}
func (muxer *Muxer) extendContext(appCtx *skipmap.StringMap[types.RoxyContext]) {
	// muxer context
	appCtx.Store("FindGroupByJid", types.FindGroupByJid(muxer.findGroupByJid))
	appCtx.Store("GetAllGroups", types.GetAllGroups(muxer.getAllGroups))
	appCtx.Store("CacheAllGroup", types.CacheAllGroup(muxer.cacheAllGroup))
	appCtx.Store("UNCacheOneGroup", types.UNCacheOneGroup(muxer.unCacheOneGroup))
	appCtx.Store("IsClientGroupAdmin", types.IsClientGroupAdmin(muxer.isClientGroupAdmin))
	appCtx.Store("IsGroupAdmin", types.IsGroupAdmin(muxer.isGroupAdmin))

	muxer.ctx = appCtx
}
