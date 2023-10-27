package command

import (
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

// UploadMediaFromUrl upload media detected from mimetype
func (runFunc *RunFuncContext) UploadMediaFromUrl(url string, vars map[string]string) (any, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadBytesMedia(bytes, vars)
}

// UploadImageFromUrl upload image from given url
func (runFunc *RunFuncContext) UploadImageFromUrl(url, caption string) (*waProto.ImageMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadImageMessageFromBytes(bytes, caption)
}

// UploadVideoFromUrl upload video from given url
func (runFunc *RunFuncContext) UploadVideoFromUrl(url, caption string) (*waProto.VideoMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadVideoMessageFromBytes(bytes, caption)
}

// UploadDocumentFromUrl upload document from given url
func (runFunc *RunFuncContext) UploadDocumentFromUrl(url, title, filename string) (*waProto.DocumentMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadDocumentMessageFromBytes(bytes, title, filename)
}

// UploadAudioFromUrl upload audio from given url
func (runFunc *RunFuncContext) UploadAudioFromUrl(url string) (*waProto.AudioMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadAudioMessageFromBytes(bytes)
}

// UploadStickerFromUrl upload sticker from given url
func (runFunc *RunFuncContext) UploadStickerFromUrl(url string) (*waProto.StickerMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return runFunc.UploadStickerMessageFromBytes(bytes)
}
