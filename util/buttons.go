package util

import (
	"github.com/itzngga/goRoxy/types"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func FixInvisibleButton(button *waProto.ButtonsMessage) *waProto.Message {
	return &waProto.Message{
		ViewOnceMessage: &waProto.FutureProofMessage{
			Message: &waProto.Message{
				MessageContextInfo: &waProto.MessageContextInfo{
					DeviceListMetadataVersion: types.Int32(3),
					DeviceListMetadata:        nil,
				},
				ButtonsMessage: button,
			},
		},
	}
}

func GenerateButton(id, cmd, text string) *waProto.ButtonsMessage_Button {
	return &waProto.ButtonsMessage_Button{
		ButtonId: types.String(CreateButtonID(id, cmd)),
		ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
			DisplayText: types.String(text),
		},
		Type: waProto.ButtonsMessage_Button_RESPONSE.Enum()}
}

func CreateTextButton(content, footer string, buttons ...*waProto.ButtonsMessage_Button) *waProto.Message {
	return FixInvisibleButton(
		&waProto.ButtonsMessage{
			HeaderType:  waProto.ButtonsMessage_TEXT.Enum(),
			ContentText: types.String(content),
			FooterText:  types.String(footer),
			Buttons:     buttons,
		},
	)
}

func CreateEmptyButton(content, footer string, buttons ...*waProto.ButtonsMessage_Button) *waProto.Message {
	return FixInvisibleButton(&waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_EMPTY.Enum(),
		ContentText: types.String(content),
		FooterText:  types.String(footer),
		Buttons:     buttons,
	},
	)
}

func CreateImageButton(content, footer string, image *waProto.ButtonsMessage_ImageMessage, buttons ...*waProto.ButtonsMessage_Button) *waProto.Message {
	return FixInvisibleButton(&waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_IMAGE.Enum(),
		ContentText: types.String(content),
		FooterText:  types.String(footer),
		Header:      image,
		Buttons:     buttons,
	},
	)
}

func CreateVideoButton(content, footer string, video *waProto.ButtonsMessage_VideoMessage, buttons ...*waProto.ButtonsMessage_Button) *waProto.Message {
	return FixInvisibleButton(&waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_VIDEO.Enum(),
		ContentText: types.String(content),
		FooterText:  types.String(footer),
		Header:      video,
		Buttons:     buttons,
	},
	)
}

func CreateLocationButton(content, footer string, location *waProto.ButtonsMessage_LocationMessage, buttons ...*waProto.ButtonsMessage_Button) *waProto.Message {
	return FixInvisibleButton(&waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_LOCATION.Enum(),
		ContentText: types.String(content),
		FooterText:  types.String(footer),
		Header:      location,
		Buttons:     buttons,
	},
	)
}

func CreateDocumentButton(content, footer string, document *waProto.ButtonsMessage_DocumentMessage, buttons ...*waProto.ButtonsMessage_Button) *waProto.Message {
	return FixInvisibleButton(&waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_DOCUMENT.Enum(),
		ContentText: types.String(content),
		FooterText:  types.String(footer),
		Header:      document,
		Buttons:     buttons,
	},
	)
}
