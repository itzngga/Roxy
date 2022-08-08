package util

import (
	"github.com/itzngga/goRoxy/helper"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

func FixInvisibleButton(button *waProto.ButtonsMessage) *waProto.Message {
	return &waProto.Message{
		ViewOnceMessage: &waProto.FutureProofMessage{
			Message: &waProto.Message{
				MessageContextInfo: &waProto.MessageContextInfo{
					DeviceListMetadataVersion: proto.Int32(3),
					DeviceListMetadata:        nil,
				},
				ButtonsMessage: button,
			},
		},
	}
}

func GenerateButton(id, cmd, text string) *waProto.Button {
	return &waProto.Button{
		ButtonId: proto.String(helper.CreateButtonID(id, cmd)),
		ButtonText: &waProto.ButtonText{
			DisplayText: proto.String(text),
		},
		Type: waProto.Button_RESPONSE.Enum()}
}

func CreateTextButton(content, footer string, buttons ...*waProto.Button) *waProto.Message {
	return FixInvisibleButton(
		&waProto.ButtonsMessage{
			HeaderType:  waProto.ButtonsMessage_TEXT.Enum(),
			ContentText: proto.String(content),
			FooterText:  proto.String(footer),
			Buttons:     buttons,
		},
	)
}

func CreateEmptyButton(content, footer string, buttons ...*waProto.Button) *waProto.Message {
	return FixInvisibleButton(&waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_EMPTY.Enum(),
		ContentText: proto.String(content),
		FooterText:  proto.String(footer),
		Buttons:     buttons,
	},
	)
}

func CreateImageButton(content, footer string, image *waProto.ButtonsMessage_ImageMessage, buttons ...*waProto.Button) *waProto.Message {
	return FixInvisibleButton(&waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_IMAGE.Enum(),
		ContentText: proto.String(content),
		FooterText:  proto.String(footer),
		Header:      image,
		Buttons:     buttons,
	},
	)
}

func CreateVideoButton(content, footer string, video *waProto.ButtonsMessage_VideoMessage, buttons ...*waProto.Button) *waProto.Message {
	return FixInvisibleButton(&waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_VIDEO.Enum(),
		ContentText: proto.String(content),
		FooterText:  proto.String(footer),
		Header:      video,
		Buttons:     buttons,
	},
	)
}

func CreateLocationButton(content, footer string, location *waProto.ButtonsMessage_LocationMessage, buttons ...*waProto.Button) *waProto.Message {
	return FixInvisibleButton(&waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_LOCATION.Enum(),
		ContentText: proto.String(content),
		FooterText:  proto.String(footer),
		Header:      location,
		Buttons:     buttons,
	},
	)
}

func CreateDocumentButton(content, footer string, document *waProto.ButtonsMessage_DocumentMessage, buttons ...*waProto.Button) *waProto.Message {
	return FixInvisibleButton(&waProto.ButtonsMessage{
		HeaderType:  waProto.ButtonsMessage_DOCUMENT.Enum(),
		ContentText: proto.String(content),
		FooterText:  proto.String(footer),
		Header:      document,
		Buttons:     buttons,
	},
	)
}
