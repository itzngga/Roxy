package roxy

import (
	"strings"

	"github.com/itzngga/Roxy/util"
	"github.com/sajari/fuzzy"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func (muxer *Muxer) GenerateSuggestionModel() {
	model := fuzzy.NewModel()
	model.SetThreshold(1)
	model.SetDepth(5)

	var words []string
	muxer.Commands.Range(func(key string, value *Command) bool {
		words = append(words, value.Name)
		words = append(words, value.Aliases...)
		return true
	})

	model.Train(words)
	muxer.SuggestionModel = model
}

func (muxer *Muxer) SuggestCommand(event *events.Message, prefix, command string) {
	suggested := muxer.SuggestionModel.Suggestions(command, false)

	if len(suggested) == 0 {
		return
	}

	for i, s := range suggested {
		suggested[i] = prefix + s
	}

	parsed := "Did you mean?: \n" + strings.Join(suggested, " or ")
	message := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String(parsed),
			ContextInfo: util.WithReply(event),
		},
	}

	muxer.SendMessage(event.Info.Chat, message)
}
