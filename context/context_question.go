package context

import (
	"encoding/json"
	"strings"

	"github.com/itzngga/Roxy/util"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

type Questions struct {
	Index        int
	Question     string
	HasReplied   bool
	Capture      bool
	CaptureMedia bool
	Reply        bool
	Answer       any
}

type QuestionState struct {
	WithEmojiReact bool
	EmojiReact     string
	Separator      string
	Ctx            *Ctx
	ActiveQuestion string
	Questions      []*Questions
	ResultChan     chan bool
}

func (q *Questions) SetAnswer(answer any) {
	switch v := answer.(type) {
	case string:
		*q.Answer.(*string) = v
	case *waE2E.Message:
		if q.CaptureMedia {
			m := util.ParseQuotedMessage(v)
			if m != nil {
				v = m
			}
			switch {
			case v.ImageMessage != nil:
				*q.Answer.(**waE2E.Message) = v
			case v.VideoMessage != nil:
				*q.Answer.(**waE2E.Message) = v
			case v.AudioMessage != nil:
				*q.Answer.(**waE2E.Message) = v
			case v.DocumentMessage != nil:
				*q.Answer.(**waE2E.Message) = v
			case v.StickerMessage != nil:
				*q.Answer.(**waE2E.Message) = v
			default:
				*q.Answer.(**waE2E.Message) = nil
			}
		} else {
			*q.Answer.(**waE2E.Message) = v
		}
	}
}

func (q *Questions) GetAnswer() string {
	switch v := q.Answer.(type) {
	case *string:
		return *q.Answer.(*string)
	case *waE2E.Message:
		result, err := json.Marshal(&v)
		if err != nil {
			return ""
		}
		return string(result)
	}
	return ""
}

// NewUserQuestion New user question engine
func (context *Ctx) NewUserQuestion() *QuestionState {
	question := &QuestionState{
		Ctx:        context,
		ResultChan: make(chan bool),
		Separator:  " | ",
	}

	return question
}

// NewUserQuestion New user question engine
func NewUserQuestion(context *Ctx) *QuestionState {
	question := &QuestionState{
		Ctx:        context,
		ResultChan: make(chan bool),
		Separator:  " | ",
	}

	return question
}

// SetParserSeparator Set parser separator for question, eg:
// /hello world | info [" | " is the separator]
func (state *QuestionState) SetParserSeparator(separator string) *QuestionState {
	state.Separator = separator
	return state
}

// SetQuestion Set a question based on question and string answer pointer
func (state *QuestionState) SetQuestion(question string, answer any) *QuestionState {
	if _, ok := answer.(*string); !ok {
		return state
	}
	state.Questions = append(state.Questions, &Questions{
		Index:    len(state.Questions) + 1,
		Question: question,
		Answer:   answer,
	})
	return state
}

// SetNoAskQuestions Set no asking question based on question and string answer pointer
func (state *QuestionState) SetNoAskQuestions(answer any) *QuestionState {
	if _, ok := answer.(*string); !ok {
		return state
	}
	state.Questions = append(state.Questions, &Questions{
		Index:    len(state.Questions) + 1,
		Question: "",
		Answer:   answer,
	})
	return state
}

// SetReplyQuestion Set a question based on message has a reply string answer pointer
func (state *QuestionState) SetReplyQuestion(question string, answer any) *QuestionState {
	if _, ok := answer.(*string); !ok {
		return state
	}
	state.Questions = append(state.Questions, &Questions{
		Index:    len(state.Questions) + 1,
		Question: question,
		Reply:    true,
		Answer:   answer,
	})
	return state
}

// SetNoAskReplyQuestion Set no asking question based on message has a reply string answer pointer
func (state *QuestionState) SetNoAskReplyQuestion(answer any) *QuestionState {
	if _, ok := answer.(*string); !ok {
		return state
	}
	state.Questions = append(state.Questions, &Questions{
		Index:    len(state.Questions) + 1,
		Question: "",
		Reply:    true,
		Answer:   answer,
	})
	return state
}

// CaptureQuestion Set a question to capture message object with json string format
func (state *QuestionState) CaptureQuestion(question string, answer **waE2E.Message) *QuestionState {
	state.Questions = append(state.Questions, &Questions{
		Index:    len(state.Questions) + 1,
		Question: question,
		Capture:  true,
		Answer:   answer,
	})
	return state
}

// NoAskCaptureQuestion Set no asking question to capture message object with json string format
func (state *QuestionState) NoAskCaptureQuestion(answer **waE2E.Message) *QuestionState {
	state.Questions = append(state.Questions, &Questions{
		Index:    len(state.Questions) + 1,
		Question: "",
		Capture:  true,
		Answer:   answer,
	})
	return state
}

// CaptureMediaQuestion Set a question to capture media object
func (state *QuestionState) CaptureMediaQuestion(question string, answer **waE2E.Message) *QuestionState {
	state.Questions = append(state.Questions, &Questions{
		Index:        len(state.Questions) + 1,
		Question:     question,
		Capture:      true,
		CaptureMedia: true,
		Answer:       answer,
	})
	return state
}

// NoAskCaptureMediaQuestion Set no asking question to capture media object
func (state *QuestionState) NoAskCaptureMediaQuestion(answer **waE2E.Message) *QuestionState {
	state.Questions = append(state.Questions, &Questions{
		Index:        len(state.Questions) + 1,
		Question:     "",
		Capture:      true,
		CaptureMedia: true,
		Answer:       answer,
	})
	return state
}

// ExecWithParser Run question engine with argument parser
func (state *QuestionState) ExecWithParser() {
	questions := strings.Split(strings.Join(state.Ctx.Arguments(), " "), state.Separator)
	if questions[0] != "" && len(state.Questions) == len(questions) {
		for i := range state.Questions {
			state.Questions[i].SetAnswer(questions[i])
		}
		return
	} else {
		state.Ctx.questionChan <- state
		defer close(state.ResultChan)

		<-state.ResultChan
		return
	}
}

// WithOkEmoji react 👌 when user answered a question
func (state *QuestionState) WithOkEmoji() *QuestionState {
	state.WithEmojiReact = true
	state.EmojiReact = "👌"
	return state
}

// WithLikeEmoji react 👍 when user answered a question
func (state *QuestionState) WithLikeEmoji() *QuestionState {
	state.WithEmojiReact = true
	state.EmojiReact = "👍"
	return state
}

// WithTimeEmoji react ⏳ when user answered a question
func (state *QuestionState) WithTimeEmoji() *QuestionState {
	state.WithEmojiReact = true
	state.EmojiReact = "⏳"
	return state
}

// WithEmoji react custom emoji when user answered a question
func (state *QuestionState) WithEmoji(emoji string) *QuestionState {
	state.WithEmojiReact = true
	state.EmojiReact = emoji
	return state
}

// Exec Run question engine without argument parser
func (state *QuestionState) Exec() {
	state.Ctx.questionChan <- state
	defer close(state.ResultChan)

	<-state.ResultChan
}
