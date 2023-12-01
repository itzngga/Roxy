package roxy

import (
	waProto "github.com/go-whatsapp/whatsmeow/binary/proto"
	"github.com/go-whatsapp/whatsmeow/types/events"
	"github.com/itzngga/Roxy/util"
	"github.com/sajari/fuzzy"
	"google.golang.org/protobuf/proto"
	"strings"
)

func (muxer *Muxer) GenerateSuggestionModel() {
	model := fuzzy.NewModel()
	model.SetThreshold(1)
	model.SetDepth(5)

	var words []string
	muxer.Commands.Range(func(key string, value *Command) bool {
		words = append(words, value.Name)
		for _, alias := range value.Aliases {
			words = append(words, alias)
		}
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

	var parsed = "Did you mean?: \n" + strings.Join(suggested, " or ")
	message := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String(parsed),
			ContextInfo: util.WithReply(event),
		},
	}

	_, _ = muxer.AppMethods.SendMessage(event.Info.Chat, message)
	return

}
