package command

import (
	"context"
	"github.com/itzngga/Roxy/util"
	"go.mau.fi/whatsmeow"
	"time"
)

type PollingOptions struct {
	Options string
	Hashed  []byte
}

type PollingState struct {
	PollId         string
	PollName       string
	PollOptions    []PollingOptions
	PollSelectable int
	PollingTimeout *time.Duration
	RunFuncCtx     *RunFuncContext
	PollingResult  []string
	ResultChan     chan bool
}

func NewPollingState(ctx *RunFuncContext) *PollingState {
	return &PollingState{
		RunFuncCtx: ctx,
		ResultChan: make(chan bool),
	}
}

func (p *PollingState) SetPollInformation(name string, options []string) *PollingState {
	var pollingOptions = make([]PollingOptions, 0)
	hashed := whatsmeow.HashPollOptions(options)
	for i, bytes := range hashed {
		pollingOptions = append(pollingOptions, PollingOptions{
			Options: options[i],
			Hashed:  bytes,
		})
	}
	p.PollName = name
	p.PollOptions = pollingOptions
	p.PollSelectable = len(options)
	return p
}

func (p *PollingState) SetTimeBasedType(timeOut time.Duration) *PollingState {
	p.PollingTimeout = &timeOut
	return p
}

func (p *PollingState) SetOnlyOnePool() *PollingState {
	p.PollSelectable = 1
	return p
}

func (p *PollingState) SetSelectableOption(count int) *PollingState {
	p.PollSelectable = count
	return p
}

func (p *PollingState) sendPollMessage() {
	var options []string
	for _, option := range p.PollOptions {
		options = append(options, option.Options)
	}

	message := p.RunFuncCtx.Client.BuildPollCreation(p.PollName, options, p.PollSelectable)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*p.RunFuncCtx.Options.SendMessageTimeout)
	defer cancel()

	response, _ := p.RunFuncCtx.Client.SendMessage(ctx, p.RunFuncCtx.MessageChat, message)
	p.PollId = response.ID
}

func (p *PollingState) Exec() []string {
	p.sendPollMessage()
	p.RunFuncCtx.PollingChan <- p
	defer close(p.ResultChan)

	_ = <-p.ResultChan
	result := util.RemoveDuplicate(p.PollingResult)
	return result
}
