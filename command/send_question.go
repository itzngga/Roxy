package command

import (
	"encoding/json"
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"strings"
)

type Questions struct {
	Index        int
	Question     string
	Capture      bool
	CaptureMedia bool
	Reply        bool
	Answer       any
}

type QuestionState struct {
	RunFuncCtx     *RunFuncContext
	ActiveQuestion string
	Questions      []*Questions
	ResultChan     chan bool
}

func (q *Questions) SetAnswer(answer any) {
	switch v := answer.(type) {
	case string:
		*q.Answer.(*string) = v
	case *waProto.Message:
		if q.CaptureMedia {
			m := util.ParseQuotedMessage(v)
			if m != nil {
				v = m
			}
			switch {
			case v.ImageMessage != nil:
				*q.Answer.(**waProto.Message) = v
			case v.VideoMessage != nil:
				*q.Answer.(**waProto.Message) = v
			case v.AudioMessage != nil:
				*q.Answer.(**waProto.Message) = v
			case v.DocumentMessage != nil:
				*q.Answer.(**waProto.Message) = v
			case v.StickerMessage != nil:
				*q.Answer.(**waProto.Message) = v
			default:
				*q.Answer.(**waProto.Message) = nil
			}
		} else {
			*q.Answer.(**waProto.Message) = v
		}
	}
}

func (q *Questions) GetAnswer() string {
	switch v := q.Answer.(type) {
	case *string:
		return *q.Answer.(*string)
	case *waProto.Message:
		result, err := json.Marshal(&v)
		if err != nil {
			return ""
		}
		return string(result)
	}
	return ""
}

// NewUserQuestion New user question engine
func NewUserQuestion(ctx *RunFuncContext) *QuestionState {
	question := &QuestionState{
		RunFuncCtx: ctx,
		ResultChan: make(chan bool),
	}

	return question
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

// CaptureQuestion Set a question to capture message object with json string format
func (state *QuestionState) CaptureQuestion(question string, answer **waProto.Message) *QuestionState {
	state.Questions = append(state.Questions, &Questions{
		Index:    len(state.Questions) + 1,
		Question: question,
		Capture:  true,
		Answer:   answer,
	})
	return state
}

func (state *QuestionState) CaptureMediaQuestion(question string, answer **waProto.Message) *QuestionState {
	state.Questions = append(state.Questions, &Questions{
		Index:        len(state.Questions) + 1,
		Question:     question,
		Capture:      true,
		CaptureMedia: true,
		Answer:       answer,
	})
	return state
}

// ExecWithParser Run question engine with argument parser
func (state *QuestionState) ExecWithParser() {
	questions := strings.Split(strings.Join(state.RunFuncCtx.Arguments, " "), " | ")
	if questions[0] != "" && len(state.Questions) == len(questions) {
		for i, _ := range state.Questions {
			state.Questions[i].SetAnswer(questions[i])
		}
		return
	} else {
		state.RunFuncCtx.QuestionChan <- state
		defer close(state.ResultChan)

		_ = <-state.ResultChan
		return
	}
}

// Exec Run question engine without argument parser
func (state *QuestionState) Exec() {
	state.RunFuncCtx.QuestionChan <- state
	defer close(state.ResultChan)

	_ = <-state.ResultChan
	return
}
