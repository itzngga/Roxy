package util

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

func FixInvisibleListMessage(listmessage *waProto.ListMessage) *waProto.Message {
	return &waProto.Message{
		ViewOnceMessage: &waProto.FutureProofMessage{
			Message: &waProto.Message{
				MessageContextInfo: &waProto.MessageContextInfo{
					DeviceListMetadataVersion: proto.Int32(2),
					DeviceListMetadata:        nil,
				},
				ListMessage: listmessage,
			},
		},
	}
}

func GenerateListMessage(title, description, buttonText, footerText string, sections ...*waProto.ListMessage_Section) *waProto.Message {
	return FixInvisibleListMessage(
		&waProto.ListMessage{
			ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
			Title:       &title,
			Description: &description,
			ButtonText:  &buttonText,
			FooterText:  &footerText,
			Sections:    sections,
		},
	)
}
func CreateSectionList(title string, rows ...*waProto.ListMessage_Row) *waProto.ListMessage_Section {
	return &waProto.ListMessage_Section{
		Title: &title,
		Rows:  rows,
	}
}

func CreateSectionRow(title, description, rowId string) *waProto.ListMessage_Row {
	return &waProto.ListMessage_Row{
		Title:       &title,
		Description: &description,
		RowId:       &rowId,
	}
}
