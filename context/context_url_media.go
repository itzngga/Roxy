package context

import (
	"github.com/itzngga/Roxy/util"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

// UploadMediaFromUrl upload media detected from mimetype
func (context *Ctx) UploadMediaFromUrl(url string, vars map[string]string) (any, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return context.UploadBytesMedia(bytes, vars)
}

// UploadImageFromUrl upload image from given url
func (context *Ctx) UploadImageFromUrl(url, caption string) (*waE2E.ImageMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return context.UploadImageMessageFromBytes(bytes, caption)
}

// UploadVideoFromUrl upload video from given url
func (context *Ctx) UploadVideoFromUrl(url, caption string) (*waE2E.VideoMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return context.UploadVideoMessageFromBytes(bytes, caption)
}

// UploadDocumentFromUrl upload document from given url
func (context *Ctx) UploadDocumentFromUrl(url, title, filename string) (*waE2E.DocumentMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return context.UploadDocumentMessageFromBytes(bytes, title, filename)
}

// UploadAudioFromUrl upload audio from given url
func (context *Ctx) UploadAudioFromUrl(url string) (*waE2E.AudioMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return context.UploadAudioMessageFromBytes(bytes)
}

// UploadStickerFromUrl upload sticker from given url
func (context *Ctx) UploadStickerFromUrl(url string) (*waE2E.StickerMessage, error) {
	bytes, err := util.DoHTTPRequest("GET", url)
	if err != nil {
		return nil, err
	}

	return context.UploadStickerMessageFromBytes(bytes)
}
