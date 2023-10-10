package core

import (
	"context"
	"fmt"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/types"
	"github.com/itzngga/Roxy/util"
	"github.com/sajari/fuzzy"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"strings"
	"time"
)

func (muxer *Muxer) GenerateSuggestionModel() {
	model := fuzzy.NewModel()
	model.SetThreshold(1)
	model.SetDepth(5)

	var words []string
	muxer.Commands.Range(func(key string, value *command.Command) bool {
		words = append(words, value.Name)
		for _, alias := range value.Aliases {
			words = append(words, alias)
		}
		return true
	})

	model.Train(words)
	muxer.SuggestionModel = model
}

func (muxer *Muxer) suggestCommand(event *events.Message, prefix, command string) {
	suggested := muxer.SuggestionModel.Suggestions(command, false)

	if len(suggested) == 0 {
		return
	}

	for i, s := range suggested {
		suggested[i] = prefix + s
	}

	var parsed = "Did you mean?: \n" + strings.Join(suggested, " or ")
	message := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        types.String(parsed),
			ContextInfo: util.WithReply(event),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var client = types.GetContext[*whatsmeow.Client](muxer.ctx, "appClient")
	_, err := client.SendMessage(ctx, event.Info.Chat, message)
	if err != nil {
		fmt.Printf("error: sending message: %v\n", err)
		return
	}

	return

}
