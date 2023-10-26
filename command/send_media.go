package command

import (
	"bufio"
	"context"
	"github.com/gabriel-vasile/mimetype"
	"github.com/itzngga/Roxy/types"
	"github.com/itzngga/Roxy/util"
	"github.com/itzngga/Roxy/util/thumbnail"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"os"
)

// UploadImageMessageFromPath upload a image from given path
func (runFunc *RunFuncContext) UploadImageMessageFromPath(path, caption string) (*waProto.ImageMessage, error) {
	imageBuff, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imageBuff.Close()

	imageInfo, _ := imageBuff.Stat()
	var imageSize = imageInfo.Size()
	imageBytes := make([]byte, imageSize)

	imageBuffer := bufio.NewReader(imageBuff)
	_, err = imageBuffer.Read(imageBytes)

	mimetypeString := mimetype.Detect(imageBytes)

	thumbnailByte := thumbnail.CreateVideoThumbnail(imageBytes)

	resp, err := runFunc.Client.Upload(context.Background(), imageBytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	thumbnail, err := runFunc.Client.Upload(context.Background(), thumbnailByte, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		thumbnailByte = nil
	}()

	return &waProto.ImageMessage{
		Caption: types.String(caption),

		Mimetype: types.String(mimetypeString.String()),

		ThumbnailDirectPath: &thumbnail.DirectPath,
		ThumbnailSha256:     thumbnail.FileSHA256,
		ThumbnailEncSha256:  thumbnail.FileEncSHA256,
		JpegThumbnail:       thumbnailByte,

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadImageMessageFromBytes upload image from given bytes
func (runFunc *RunFuncContext) UploadImageMessageFromBytes(bytes []byte, caption string) (*waProto.ImageMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	thumbnailByte := thumbnail.CreateImageThumbnail(bytes)

	resp, err := runFunc.Client.Upload(context.Background(), bytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	thumbnail, err := runFunc.Client.Upload(context.Background(), thumbnailByte, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		thumbnailByte = nil
		bytes = nil
	}()

	return &waProto.ImageMessage{
		ContextInfo: util.WithReply(runFunc.MessageEvent),

		Caption:  types.String(caption),
		Mimetype: types.String(mimetypeString.String()),

		ThumbnailDirectPath: &thumbnail.DirectPath,
		ThumbnailSha256:     thumbnail.FileSHA256,
		ThumbnailEncSha256:  thumbnail.FileEncSHA256,
		JpegThumbnail:       thumbnailByte,

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadVideoMessageFromPath upload video from given path
func (runFunc *RunFuncContext) UploadVideoMessageFromPath(path, caption string) (*waProto.VideoMessage, error) {
	videoBuff, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer videoBuff.Close()

	videoInfo, _ := videoBuff.Stat()
	var videoSize = videoInfo.Size()
	videoBytes := make([]byte, videoSize)

	videoBuffer := bufio.NewReader(videoBuff)
	_, err = videoBuffer.Read(videoBytes)

	mimetypeString := mimetype.Detect(videoBytes)

	thumbnailByte := thumbnail.CreateVideoThumbnail(videoBytes)

	resp, err := runFunc.Client.Upload(context.Background(), videoBytes, whatsmeow.MediaVideo)
	if err != nil {
		return nil, err
	}

	thumbnail, err := runFunc.Client.Upload(context.Background(), thumbnailByte, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		thumbnailByte = nil
	}()

	return &waProto.VideoMessage{
		Caption:  types.String(caption),
		Mimetype: types.String(mimetypeString.String()),

		ThumbnailDirectPath: &thumbnail.DirectPath,
		ThumbnailSha256:     thumbnail.FileSHA256,
		ThumbnailEncSha256:  thumbnail.FileEncSHA256,
		JpegThumbnail:       thumbnailByte,

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadVideoMessageFromBytes upload video from given bytes
func (runFunc *RunFuncContext) UploadVideoMessageFromBytes(bytes []byte, caption string) (*waProto.VideoMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	thumbnailByte := thumbnail.CreateVideoThumbnail(bytes)

	resp, err := runFunc.Client.Upload(context.Background(), bytes, whatsmeow.MediaVideo)
	if err != nil {
		return nil, err
	}

	thumbnail, err := runFunc.Client.Upload(context.Background(), thumbnailByte, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		thumbnailByte = nil
		bytes = nil
	}()

	return &waProto.VideoMessage{
		ContextInfo: util.WithReply(runFunc.MessageEvent),

		Caption:  types.String(caption),
		Mimetype: types.String(mimetypeString.String()),

		ThumbnailDirectPath: &thumbnail.DirectPath,
		ThumbnailSha256:     thumbnail.FileSHA256,
		ThumbnailEncSha256:  thumbnail.FileEncSHA256,
		JpegThumbnail:       thumbnailByte,

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadStickerMessageFromPath upload sticker from given path
func (runFunc *RunFuncContext) UploadStickerMessageFromPath(path string) (*waProto.StickerMessage, error) {
	stickerBuff, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer stickerBuff.Close()

	stickerInfo, _ := stickerBuff.Stat()
	var stickerSize = stickerInfo.Size()
	stickerBytes := make([]byte, stickerSize)

	stickerBuffer := bufio.NewReader(stickerBuff)
	_, err = stickerBuffer.Read(stickerBytes)

	resp, err := runFunc.Client.Upload(context.Background(), stickerBytes, whatsmeow.MediaImage)

	if err != nil {
		return nil, err
	}

	return &waProto.StickerMessage{
		Mimetype: types.String("image/webp"),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadStickerMessageFromBytes upload sticker from given bytes
func (runFunc *RunFuncContext) UploadStickerMessageFromBytes(bytes []byte) (*waProto.StickerMessage, error) {
	resp, err := runFunc.Client.Upload(context.Background(), bytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		bytes = nil
	}()

	return &waProto.StickerMessage{
		ContextInfo: util.WithReply(runFunc.MessageEvent),

		Mimetype: types.String("image/webp"),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadDocumentMessageFromPath upload document from given path
func (runFunc *RunFuncContext) UploadDocumentMessageFromPath(path, title string) (*waProto.DocumentMessage, error) {
	documentBuff, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer documentBuff.Close()

	documentInfo, _ := documentBuff.Stat()
	var documentSize = documentInfo.Size()
	documentBytes := make([]byte, documentSize)

	documentBuffer := bufio.NewReader(documentBuff)
	_, err = documentBuffer.Read(documentBytes)

	mimetypeString := mimetype.Detect(documentBytes)

	resp, err := runFunc.Client.Upload(context.Background(), documentBytes, whatsmeow.MediaDocument)
	if err != nil {
		return nil, err
	}

	return &waProto.DocumentMessage{
		Title:    types.String(title),
		FileName: types.String(documentInfo.Name()),
		Mimetype: types.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadDocumentMessageFromBytes upload document from given bytes
func (runFunc *RunFuncContext) UploadDocumentMessageFromBytes(bytes []byte, title, filename string) (*waProto.DocumentMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	resp, err := runFunc.Client.Upload(context.Background(), bytes, whatsmeow.MediaDocument)

	if err != nil {
		return nil, err
	}

	defer func() {
		bytes = nil
	}()

	return &waProto.DocumentMessage{
		Title:    types.String(title),
		FileName: types.String(filename),
		Mimetype: types.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadAudioMessageFromPath upload audio message from given path
func (runFunc *RunFuncContext) UploadAudioMessageFromPath(path string) (*waProto.AudioMessage, error) {
	audioBuff, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer audioBuff.Close()

	audioInfo, _ := audioBuff.Stat()
	var audioSize = audioInfo.Size()
	audioBytes := make([]byte, audioSize)

	audioBuffer := bufio.NewReader(audioBuff)
	_, err = audioBuffer.Read(audioBytes)

	mimetypeString := mimetype.Detect(audioBytes)

	resp, err := runFunc.Client.Upload(context.Background(), audioBytes, whatsmeow.MediaAudio)

	if err != nil {
		return nil, err
	}

	return &waProto.AudioMessage{
		Mimetype: types.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadAudioMessageFromBytes upload audio from given bytes
func (runFunc *RunFuncContext) UploadAudioMessageFromBytes(bytes []byte) (*waProto.AudioMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	resp, err := runFunc.Client.Upload(context.Background(), bytes, whatsmeow.MediaAudio)
	if err != nil {
		return nil, err
	}

	defer func() {
		bytes = nil
	}()

	return &waProto.AudioMessage{
		Mimetype: types.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}
