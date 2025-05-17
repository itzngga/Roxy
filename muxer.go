package roxy

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/itzngga/Roxy/context"
	"github.com/itzngga/Roxy/options"
	"github.com/itzngga/Roxy/types"
	"github.com/itzngga/Roxy/util"
	"github.com/puzpuzpuz/xsync"
	"github.com/sajari/fuzzy"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type Muxer struct {
	Options              *options.Options                             `json:"options,omitempty"`
	Log                  waLog.Logger                                 `json:"log,omitempty"`
	MessageTimeout       time.Duration                                `json:"message_timeout,omitempty"`
	Categories           *xsync.MapOf[string, string]                 `json:"categories,omitempty"`
	GlobalMiddlewares    *xsync.MapOf[string, context.MiddlewareFunc] `json:"global_middlewares,omitempty"`
	Middlewares          *xsync.MapOf[string, context.MiddlewareFunc] `json:"middlewares,omitempty"`
	Commands             *xsync.MapOf[string, *Command]               `json:"commands,omitempty"`
	CommandResponseCache *xsync.MapOf[string, *waProto.Message]       `json:"command_response_cache,omitempty"`
	QuestionState        *xsync.MapOf[string, *context.QuestionState] `json:"question_state,omitempty"`
	PollingState         *xsync.MapOf[string, *context.PollingState]  `json:"polling_state,omitempty"`
	GroupCache           *xsync.MapOf[string, []*waTypes.GroupInfo]   `json:"group_cache,omitempty"`
	Locals               *xsync.MapOf[string, string]                 `json:"locals,omitempty"`

	QuestionChan    chan *context.QuestionState                           `json:"question_chan,omitempty"`
	PollingChan     chan *context.PollingState                            `json:"polling_chan,omitempty"`
	SuggestionModel *fuzzy.Model                                          `json:"suggestion_model,omitempty"`
	CommandParser   func(str string) (prefix string, cmd string, ok bool) `json:"command_parser,omitempty"`

	types.AppMethods
}

func NewMuxer(log waLog.Logger, options *options.Options, appMethods types.AppMethods) *Muxer {
	muxer := &Muxer{
		Locals:               xsync.NewMapOf[string](),
		Commands:             xsync.NewMapOf[*Command](),
		GlobalMiddlewares:    xsync.NewMapOf[context.MiddlewareFunc](),
		Middlewares:          xsync.NewMapOf[context.MiddlewareFunc](),
		CommandResponseCache: xsync.NewMapOf[*waProto.Message](),
		QuestionState:        xsync.NewMapOf[*context.QuestionState](),
		PollingState:         xsync.NewMapOf[*context.PollingState](),
		Categories:           xsync.NewMapOf[string](),
		GroupCache:           xsync.NewMapOf[[]*waTypes.GroupInfo](),
		QuestionChan:         make(chan *context.QuestionState),
		PollingChan:          make(chan *context.PollingState),
		MessageTimeout:       options.SendMessageTimeout,
		Options:              options,
		Log:                  log,
		CommandParser:        util.ParseCmd,
		AppMethods:           appMethods,
	}

	go muxer.handleQuestionStateChannel()
	go muxer.handlePollingStateChannel()

	return muxer
}

func (muxer *Muxer) AddCommandParser(context func(str string) (prefix string, cmd string, ok bool)) {
	muxer.CommandParser = context
}

func (muxer *Muxer) Clean() {
	muxer.Categories.Range(func(key string, category string) bool {
		muxer.Categories.Delete(key)
		return true
	})
	muxer.Commands.Range(func(key string, cmd *Command) bool {
		muxer.Commands.Delete(key)
		return true
	})
	muxer.Middlewares.Range(func(key string, middleware context.MiddlewareFunc) bool {
		muxer.Middlewares.Delete(key)
		return true
	})
	muxer.Locals.Range(func(key string, middleware string) bool {
		muxer.Locals.Delete(key)
		return true
	})
}

func (muxer *Muxer) handlePollingStateChannel() {
	for message := range muxer.PollingChan {
		muxer.PollingState.Store(message.PollId, message)
		if message.PollingTimeout != nil {
			go func() {
				timeout := time.NewTimer(*message.PollingTimeout)
				<-timeout.C
				message.ResultChan <- true
				timeout.Stop()
				muxer.PollingState.Delete(message.PollId)
			}()
		} else {
			go func() {
				timeout := time.NewTimer(time.Minute * 10)
				<-timeout.C
				message.ResultChan <- true
				timeout.Stop()
				muxer.PollingState.Delete(message.PollId)
			}()
		}
	}
}

func (muxer *Muxer) handleQuestionStateChannel() {
	for message := range muxer.QuestionChan {
		muxer.QuestionState.Delete(message.Ctx.Number())
		for _, question := range message.Questions {
			if question.GetAnswer() == "" {
				message.ActiveQuestion = question.Question
				muxer.QuestionState.Store(message.Ctx.Number(), message)
				if question.Question != "" {
					message.Ctx.SendReplyMessage(question.Question)
				}
				break
			}
		}
	}
}

func (muxer *Muxer) addEmbedCommands() {
	categories := Categories.Get()
	for _, cat := range categories {
		muxer.Categories.Store(cat, cat)
	}
	commands := Commands.Get()
	for _, cmd := range commands {
		muxer.AddCommand(cmd)
	}
	middlewares := Middlewares.Get()
	for _, mid := range middlewares {
		muxer.AddMiddleware(mid)
	}
	globalMiddleware := GlobalMiddlewares.Get()
	for _, mid := range globalMiddleware {
		muxer.AddGlobalMiddleware(mid)
	}

	if muxer.Options.CommandSuggestion {
		muxer.GenerateSuggestionModel()
	}
}

func (muxer *Muxer) AddGlobalMiddleware(middleware context.MiddlewareFunc) {
	muxer.GlobalMiddlewares.Store(uuid.New().String(), middleware)
}

func (muxer *Muxer) AddMiddleware(middleware context.MiddlewareFunc) {
	muxer.Middlewares.Store(uuid.New().String(), middleware)
}

func (muxer *Muxer) AddCommand(cmd *Command) {
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

func (muxer *Muxer) GetActiveCommand() []*Command {
	cmd := make([]*Command, 0)
	muxer.Commands.Range(func(key string, value *Command) bool {
		// filter alias commands
		if key == value.Name {
			cmd = append(cmd, value)
		}
		return true
	})

	return cmd
}

func (muxer *Muxer) GetActiveGlobalMiddleware() []context.MiddlewareFunc {
	middleware := make([]context.MiddlewareFunc, 0)
	muxer.GlobalMiddlewares.Range(func(key string, value context.MiddlewareFunc) bool {
		middleware = append(middleware, value)
		return true
	})

	return middleware
}

func (muxer *Muxer) GetActiveMiddleware() []context.MiddlewareFunc {
	middleware := make([]context.MiddlewareFunc, 0)
	muxer.Middlewares.Range(func(key string, value context.MiddlewareFunc) bool {
		middleware = append(middleware, value)
		return true
	})

	return middleware
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
	timeout := time.NewTimer(muxer.Options.CommandResponseCacheTimeout)
	<-timeout.C
	muxer.CommandResponseCache.Delete(cmd)
	timeout.Stop()
}

func (muxer *Muxer) globalMiddlewareProcessing(c *whatsmeow.Client, evt *events.Message) bool {
	if muxer.GlobalMiddlewares.Size() >= 1 {
		ctx := context.NewCtx(muxer.Locals)
		ctx.SetClient(c)
		ctx.SetLogger(muxer.Log)
		ctx.SetOptions(muxer.Options)
		ctx.SetMessageEvent(evt)
		ctx.SetClientJID(muxer.ClientJID())
		ctx.SetClientMethods(muxer)
		ctx.SetQuestionChan(muxer.QuestionChan)
		ctx.SetPollingChan(muxer.PollingChan)
		defer context.ReleaseCtx(ctx)

		midAreOk := true
		muxer.GlobalMiddlewares.Range(func(key string, value context.MiddlewareFunc) bool {
			if !value(ctx) {
				midAreOk = false
				return false
			}
			return true
		})
		return midAreOk
	}

	return true
}

func (muxer *Muxer) handlePollingState(c *whatsmeow.Client, evt *events.Message) {
	if evt.Message.PollUpdateMessage.PollCreationMessageKey == nil && evt.Message.PollUpdateMessage.PollCreationMessageKey.ID == nil {
		return
	}

	pollingState, ok := muxer.PollingState.Load(*evt.Message.PollUpdateMessage.PollCreationMessageKey.ID)
	if ok {
		pollMessage, err := c.DecryptPollVote(evt)
		if err != nil {
			return
		}

		var result []string
		for _, selectedOption := range pollMessage.SelectedOptions {
			for _, option := range pollingState.PollOptions {
				if bytes.Equal(selectedOption, option.Hashed) {
					result = append(result, option.Options)
					break
				}
			}
		}

		pollingState.PollingResult = append(pollingState.PollingResult, result...)
		if pollingState.PollingTimeout == nil {
			pollingState.ResultChan <- true
			muxer.PollingState.Delete(pollingState.PollId)
		}
	}
}

func (muxer *Muxer) handleQuestionState(c *whatsmeow.Client, evt *events.Message, number, parsedMsg string) {
	questionState, _ := muxer.QuestionState.Load(number)
	if strings.Contains(parsedMsg, "cancel") || strings.Contains(parsedMsg, "batal") {
		muxer.QuestionState.Delete(number)
		return
	} else {
		if questionState.WithEmojiReact {
			muxer.SendEmojiMessage(evt, questionState.EmojiReact)
		}

		jids := []waTypes.MessageID{
			evt.Info.ID,
		}
		c.MarkRead(jids, evt.Info.Timestamp, evt.Info.Chat, evt.Info.Sender)

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
				continue
			} else if question.Question != questionState.ActiveQuestion && question.GetAnswer() == "" {
				questionState.ActiveQuestion = question.Question
				if question.Question != "" {
					questionState.Ctx.SendReplyMessage(question.Question)
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
	if !muxer.isAllowedSource(evt) || !muxer.isFromValidSender(evt) || evt.Info.ID == "status@broadcast" {
		return
	}

	if poll := evt.Message.GetPollUpdateMessage(); poll != nil {
		muxer.handlePollingState(c, evt)
		return
	}

	number := evt.Info.Sender.ToNonAD().String()
	parsed := util.ParseMessageText(evt)

	if _, ok := muxer.QuestionState.Load(number); ok {
		muxer.handleQuestionState(c, evt, number, parsed)
		return
	}

	if !muxer.globalMiddlewareProcessing(c, evt) {
		return
	}

	prefix, cmd, isCmd := muxer.CommandParser(parsed)
	cmdLoad, isAvailable := muxer.Commands.Load(cmd)

	if isCmd && muxer.Options.CommandSuggestion && !isAvailable {
		muxer.SuggestCommand(evt, prefix, cmd)
		return
	}

	if isCmd && isAvailable {
		ctx := context.NewCtx(muxer.Locals)
		defer context.ReleaseCtx(ctx)

		ctx.SetClient(c)
		ctx.SetLogger(muxer.Log)
		ctx.SetOptions(muxer.Options)
		ctx.SetMessageEvent(evt)
		ctx.SetClientJID(muxer.ClientJID())
		ctx.SetParsedMsg(parsed)
		ctx.SetPrefix(prefix)
		ctx.SetClientMethods(muxer)
		ctx.SetQuestionChan(muxer.QuestionChan)
		ctx.SetPollingChan(muxer.PollingChan)

		if len(cmdLoad.SubCommands) > 0 && len(ctx.Arguments) > 0 {
			if child, ok := cmdLoad.SubCommands[ctx.Arguments[0]]; ok {
				ctx.Arguments = ctx.Arguments[1:]
				go muxer.markAsReadAndLogCommand(c, evt, number, fmt.Sprintf("%s+%s", cmdLoad.Name, child.Name))
				muxer.executeCommand(evt, ctx, child, parsed)
				return
			}
		}

		go muxer.markAsReadAndLogCommand(c, evt, number, cmdLoad.Name)
		muxer.executeCommand(evt, ctx, cmdLoad, parsed)
	}
}

func (muxer *Muxer) executeCommand(evt *events.Message, ctx *context.Ctx, cmdLoad *Command, parsed string) {
	if !muxer.passesCommandGuards(evt, cmdLoad) {
		return
	}

	if !muxer.runMiddlewares(ctx, cmdLoad) {
		return
	}

	var msg *waProto.Message
	if cmdLoad.Cache {
		msg = muxer.getCachedCommandResponse(parsed)
		if msg == nil {
			msg = cmdLoad.RunFunc(ctx)
		}
	} else {
		msg = cmdLoad.RunFunc(ctx)
	}

	if msg != nil {
		if _, err := muxer.SendMessage(evt.Info.Chat, msg); err == nil && cmdLoad.Cache {
			muxer.setCacheCommandResponse(parsed, msg)
		}
	}
}

func (muxer *Muxer) isAllowedSource(evt *events.Message) bool {
	switch {
	case muxer.Options.AllowFromPrivate && !muxer.Options.AllowFromGroup && evt.Info.IsGroup:
		return false
	case muxer.Options.AllowFromGroup && !muxer.Options.AllowFromPrivate && !evt.Info.IsGroup:
		return false
	case !muxer.Options.AllowFromGroup && !muxer.Options.AllowFromPrivate:
		return false
	}
	return true
}

func (muxer *Muxer) isFromValidSender(evt *events.Message) bool {
	return !muxer.Options.OnlyFromSelf || evt.Info.IsFromMe
}

func (muxer *Muxer) passesCommandGuards(evt *events.Message, cmd *Command) bool {
	if cmd.GroupOnly && !evt.Info.IsGroup {
		return false
	}
	if cmd.PrivateOnly && evt.Info.IsGroup {
		return false
	}
	if cmd.OnlyAdminGroup && evt.Info.IsGroup {
		if ok, _ := muxer.IsGroupAdmin(evt.Info.Chat, evt.Info.Sender); !ok {
			return false
		}
	}
	if cmd.OnlyIfBotAdmin && evt.Info.IsGroup {
		if ok, _ := muxer.IsClientGroupAdmin(evt.Info.Chat); !ok {
			return false
		}
	}
	return true
}

func (muxer *Muxer) markAsReadAndLogCommand(c *whatsmeow.Client, evt *events.Message, number string, cmd string) {
	if muxer.Options.WithCommandLog {
		muxer.Log.Infof("[CMD] [%s] command > %s", number, cmd)
	}
	if err := c.MarkRead([]waTypes.MessageID{evt.Info.ID}, evt.Info.Timestamp, evt.Info.Chat, evt.Info.Sender); err != nil {
		muxer.Log.Errorf("read message error: %v", err)
	}
}

func (muxer *Muxer) runMiddlewares(ctx *context.Ctx, cmd *Command) bool {
	if muxer.Middlewares.Size() > 0 {
		allPassed := true
		muxer.Middlewares.Range(func(_ string, middleware context.MiddlewareFunc) bool {
			if !middleware(ctx) {
				allPassed = false
				return false
			}
			return true
		})
		if !allPassed {
			return false
		}
	}

	if cmd.Middleware != nil && !cmd.Middleware(ctx) {
		return false
	}

	return true
}

func (muxer *Muxer) SendEmojiMessage(event *events.Message, emoji string) {
	id := event.Info.ID
	chat := event.Info.Chat
	sender := event.Info.Sender
	key := &waCommon.MessageKey{
		FromMe:    proto.Bool(true),
		ID:        proto.String(id),
		RemoteJID: proto.String(chat.String()),
	}

	if !sender.IsEmpty() && sender.User != muxer.AppMethods.ClientJID().ToNonAD().String() {
		key.FromMe = proto.Bool(false)
		key.Participant = proto.String(sender.ToNonAD().String())
	}

	message := &waProto.Message{
		ReactionMessage: &waProto.ReactionMessage{
			Key:               key,
			Text:              proto.String(emoji),
			SenderTimestampMS: proto.Int64(time.Now().UnixMilli()),
		},
	}

	muxer.SendMessage(event.Info.Chat, message)
}
