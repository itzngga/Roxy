package command

import (
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func (runFunc *RunFuncContext) UploadImageFromUrl(url, caption string) (*waProto.ImageMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadImageMessageFromBytes(bytes, caption)
}

func (runFunc *RunFuncContext) UploadVideoFromUrl(url, caption string) (*waProto.VideoMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadVideoMessageFromBytes(bytes, caption)
}

func (runFunc *RunFuncContext) UploadDocumentFromUrl(url, title, filename string) (*waProto.DocumentMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadDocumentMessageFromBytes(bytes, title, filename)
}

func (runFunc *RunFuncContext) UploadAudioFromUrl(url string) (*waProto.AudioMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadAudioMessageFromBytes(bytes)
}

func (runFunc *RunFuncContext) UploadStickerFromUrl(url string) (*waProto.StickerMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadStickerMessageFromBytes(bytes)
}
