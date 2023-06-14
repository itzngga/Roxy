package command

import (
	"encoding/json"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"strings"
)

func (q *Questions) SetAnswer(answer any) {
	switch v := answer.(type) {
	case string:
		*q.Answer.(*string) = v
	case *waProto.Message:
		result, _ := json.Marshal(&v)
		*q.Answer.(*string) = string(result)
	}
}

func (q *Questions) GetAnswer() string {
	return *q.Answer.(*string)
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

// CaptureQuestion Set a question to capture message object with json string format
func (state *QuestionState) CaptureQuestion(question string, answer any) *QuestionState {
	if _, ok := answer.(*string); !ok {
		return state
	}

	state.Questions = append(state.Questions, &Questions{
		Index:    len(state.Questions) + 1,
		Question: question,
		Capture:  true,
		Answer:   answer,
	})
	return state
}

// Exec Run question engine process
func (state *QuestionState) Exec() {
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
