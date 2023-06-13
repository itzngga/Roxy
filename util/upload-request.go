package util

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

// @ Documentation fo gofast
// @ https://github.com/cloudingcity/gofast

func UploadImageFromUrl(c *whatsmeow.Client, m *events.Message, url, caption string) (*waProto.ImageMessage, error) {
	bytes, err := DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return UploadImageMessageFromBytes(c, m, bytes, caption)
}

func UploadVideoFromUrl(c *whatsmeow.Client, m *events.Message, url, caption string) (*waProto.VideoMessage, error) {
	bytes, err := DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return UploadVideoMessageFromBytes(c, m, bytes, caption)
}

func UploadDocumentFromUrl(c *whatsmeow.Client, url, title, filename string) (*waProto.DocumentMessage, error) {
	bytes, err := DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return UploadDocumentMessageFromBytes(c, bytes, title, filename)
}

func UploadAudioFromUrl(c *whatsmeow.Client, url string) (*waProto.AudioMessage, error) {
	bytes, err := DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return UploadAudioMessageFromBytes(c, bytes)
}

func UploadStickerFromUrl(c *whatsmeow.Client, m *events.Message, url string) (*waProto.StickerMessage, error) {
	bytes, err := DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return UploadStickerMessageFromBytes(c, m, bytes)
}
