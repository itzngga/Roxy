package util

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

func FixInvisibleTemplate(template *waProto.TemplateMessage) *waProto.Message {
	return &waProto.Message{
		ViewOnceMessage: &waProto.FutureProofMessage{
			Message: &waProto.Message{
				MessageContextInfo: &waProto.MessageContextInfo{
					DeviceListMetadataVersion: proto.Int32(2),
					DeviceListMetadata:        nil,
				},
				TemplateMessage: template,
			},
		},
	}
}

func CreateHydratedTemplateButton(content, footer string, buttons ...*waProto.HydratedTemplateButton) *waProto.Message {
	return FixInvisibleTemplate(&waProto.TemplateMessage{
		HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
			HydratedContentText: proto.String(content),
			HydratedFooterText:  proto.String(footer),
			HydratedButtons:     buttons,
		},
	})
}

func CreateHydratedTemplateImageButton(content, footer string, image *waProto.ImageMessage, buttons ...*waProto.HydratedTemplateButton) *waProto.Message {
	return FixInvisibleTemplate(&waProto.TemplateMessage{
		HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
			HydratedContentText: proto.String(content),
			HydratedFooterText:  proto.String(footer),
			HydratedButtons:     buttons,
			Title: &waProto.TemplateMessage_HydratedFourRowTemplate_ImageMessage{
				ImageMessage: image,
			},
		},
	})
}

func CreateHydratedTemplateVideoButton(content, footer string, video *waProto.VideoMessage, buttons ...*waProto.HydratedTemplateButton) *waProto.Message {
	return FixInvisibleTemplate(&waProto.TemplateMessage{
		HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
			HydratedContentText: proto.String(content),
			HydratedFooterText:  proto.String(footer),
			HydratedButtons:     buttons,
			Title: &waProto.TemplateMessage_HydratedFourRowTemplate_VideoMessage{
				VideoMessage: video,
			},
		},
	})
}

func CreateHydratedTemplateDocumentButton(content, footer string, document *waProto.DocumentMessage, buttons ...*waProto.HydratedTemplateButton) *waProto.Message {
	return FixInvisibleTemplate(&waProto.TemplateMessage{
		HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
			HydratedContentText: proto.String(content),
			HydratedFooterText:  proto.String(footer),
			HydratedButtons:     buttons,
			Title: &waProto.TemplateMessage_HydratedFourRowTemplate_DocumentMessage{
				DocumentMessage: document,
			},
		},
	})
}

func CreateHydratedTemplateLocationButton(content, footer string, location *waProto.LocationMessage, buttons ...*waProto.HydratedTemplateButton) *waProto.Message {
	return FixInvisibleTemplate(&waProto.TemplateMessage{
		HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
			HydratedContentText: proto.String(content),
			HydratedFooterText:  proto.String(footer),
			HydratedButtons:     buttons,
			Title: &waProto.TemplateMessage_HydratedFourRowTemplate_LocationMessage{
				LocationMessage: location,
			},
		},
	})
}

func CreateHydratedTemplateHydratedTitleButton(content, footer, title string, buttons ...*waProto.HydratedTemplateButton) *waProto.Message {
	return FixInvisibleTemplate(&waProto.TemplateMessage{
		HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
			HydratedContentText: proto.String(content),
			HydratedFooterText:  proto.String(footer),
			HydratedButtons:     buttons,
			Title: &waProto.TemplateMessage_HydratedFourRowTemplate_HydratedTitleText{
				HydratedTitleText: title,
			},
		},
	})
}

func GenerateHydratedUrlButton(text, url string) *waProto.HydratedTemplateButton {
	return &waProto.HydratedTemplateButton{
		HydratedButton: &waProto.HydratedTemplateButton_UrlButton{
			UrlButton: &waProto.HydratedTemplateButton_HydratedURLButton{
				DisplayText: proto.String(text),
				Url:         proto.String(url),
			},
		},
	}
}

func GenerateHydratedCallButton(text, number string) *waProto.HydratedTemplateButton {
	return &waProto.HydratedTemplateButton{
		HydratedButton: &waProto.HydratedTemplateButton_CallButton{
			CallButton: &waProto.HydratedTemplateButton_HydratedCallButton{
				DisplayText: proto.String(text),
				PhoneNumber: proto.String(number),
			},
		},
	}
}

func GenerateHydratedQuickReplyButton(text, id, cmd string) *waProto.HydratedTemplateButton {
	return &waProto.HydratedTemplateButton{
		HydratedButton: &waProto.HydratedTemplateButton_QuickReplyButton{
			QuickReplyButton: &waProto.HydratedTemplateButton_HydratedQuickReplyButton{
				DisplayText: proto.String(text),
				Id:          proto.String(CreateButtonID(id, cmd)),
			},
		},
	}
}
