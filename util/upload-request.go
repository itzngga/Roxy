package util

import (
	"github.com/itzngga/goRoxy/util/gofast"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

// @ Documentation fo gofast
// @ https://github.com/cloudingcity/gofast

func UploadImageFromUrl(c *whatsmeow.Client, url, caption string) (*waProto.ImageMessage, error) {
	fast := gofast.New(gofast.Config{
		ResponseDecoder: gofast.ByteDecoder,
	})
	var bytes []byte
	if err := fast.Get(url, &bytes, nil); err != nil {
		return nil, err
	}

	return UploadImageMessageFromBytes(c, bytes, caption)
}

func UploadVideoFromUrl(c *whatsmeow.Client, url, caption string) (*waProto.VideoMessage, error) {
	fast := gofast.New(gofast.Config{
		ResponseDecoder: gofast.ByteDecoder,
	})
	var bytes []byte
	if err := fast.Get(url, &bytes, nil); err != nil {
		return nil, err
	}

	return UploadVideoMessageFromBytes(c, bytes, caption)
}

func UploadDocumentFromUrl(c *whatsmeow.Client, url, title, filename string) (*waProto.DocumentMessage, error) {
	fast := gofast.New(gofast.Config{
		ResponseDecoder: gofast.ByteDecoder,
	})
	var bytes []byte
	if err := fast.Get(url, &bytes, nil); err != nil {
		return nil, err
	}

	return UploadDocumentMessageFromBytes(c, bytes, title, filename)
}

func UploadAudioFromUrl(c *whatsmeow.Client, url string) (*waProto.AudioMessage, error) {
	fast := gofast.New(gofast.Config{
		ResponseDecoder: gofast.ByteDecoder,
	})
	var bytes []byte
	if err := fast.Get(url, &bytes, nil); err != nil {
		return nil, err
	}

	return UploadAudioMessageFromBytes(c, bytes)
}

func UploadStickerFromUrl(c *whatsmeow.Client, url string) (*waProto.StickerMessage, error) {
	fast := gofast.New(gofast.Config{
		ResponseDecoder: gofast.ByteDecoder,
	})
	var bytes []byte
	if err := fast.Get(url, &bytes, nil); err != nil {
		return nil, err
	}

	return UploadStickerMessageFromBytes(c, bytes)
}
