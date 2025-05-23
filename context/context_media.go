package context

import (
	"bufio"
	"errors"
	"os"
	"slices"

	"github.com/gabriel-vasile/mimetype"
	"github.com/itzngga/Roxy/util"
	"github.com/itzngga/Roxy/util/thumbnail"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"

	context2 "context"
)

var (
	imageExtensions = []string{".png", ".jpg", ".webp", ".gif", ".bmp", ".ico", ".svg"}
	videoExtensions = []string{".mp4", ".mov", ".mpeg", ".webp", ".3gp", ".avi", ".mkv"}
	audioExtensions = []string{".mp3", ".ogg", ".m4a", ".wav", ".flac"}
)

// UploadBytesMedia upload bytes media based from mimetype
func (context *Ctx) UploadBytesMedia(bytes []byte, vars map[string]string) (any, error) {
	mimetypeString := mimetype.Detect(bytes)
	extension := mimetypeString.Extension()

	for _, imageExtension := range imageExtensions {
		if extension == imageExtension {
			caption, ok := vars["caption"]
			if !ok {
				return nil, errors.New("error: missing caption in vars map")
			}

			return context.UploadImageMessageFromBytes(bytes, caption)
		}
	}

	for _, videoExtension := range videoExtensions {
		if extension == videoExtension {
			caption, ok := vars["caption"]
			if !ok {
				return nil, errors.New("error: missing caption in vars map")
			}

			return context.UploadVideoMessageFromBytes(bytes, caption)
		}
	}

	if slices.Contains(audioExtensions, extension) {
		return context.UploadAudioMessageFromBytes(bytes)
	}

	title, ok := vars["title"]
	if !ok {
		return nil, errors.New("error: missing title in vars map")
	}

	filename, ok := vars["filename"]
	if !ok {
		return nil, errors.New("error: missing filename in vars map")
	}

	return context.UploadDocumentMessageFromBytes(bytes, title, filename)
}

// UploadImageMessageFromPath upload a image from given path
func (context *Ctx) UploadImageMessageFromPath(path, caption string) (*waE2E.ImageMessage, error) {
	imageBuff, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imageBuff.Close()

	imageInfo, _ := imageBuff.Stat()
	imageSize := imageInfo.Size()
	imageBytes := make([]byte, imageSize)

	imageBuffer := bufio.NewReader(imageBuff)
	_, err = imageBuffer.Read(imageBytes)
	if err != nil {
		return nil, err
	}

	mimetypeString := mimetype.Detect(imageBytes)

	thumbnailByte := thumbnail.CreateVideoThumbnail(imageBytes)

	resp, err := context.Client().Upload(context2.Background(), imageBytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	thumbnail, err := context.Client().Upload(context2.Background(), thumbnailByte, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		thumbnailByte = nil
	}()

	return &waE2E.ImageMessage{
		Caption: proto.String(caption),

		Mimetype: proto.String(mimetypeString.String()),

		ThumbnailDirectPath: &thumbnail.DirectPath,
		ThumbnailSHA256:     thumbnail.FileSHA256,
		ThumbnailEncSHA256:  thumbnail.FileEncSHA256,
		JPEGThumbnail:       thumbnailByte,

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadImageMessageFromBytes upload image from given bytes
func (context *Ctx) UploadImageMessageFromBytes(bytes []byte, caption string) (*waE2E.ImageMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	thumbnailByte := thumbnail.CreateImageThumbnail(bytes)

	resp, err := context.Client().Upload(context2.Background(), bytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	thumbnail, err := context.Client().Upload(context2.Background(), thumbnailByte, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		thumbnailByte = nil
		bytes = nil
	}()

	return &waE2E.ImageMessage{
		ContextInfo: util.WithReply(context.MessageEvent()),

		Caption:  proto.String(caption),
		Mimetype: proto.String(mimetypeString.String()),

		ThumbnailDirectPath: &thumbnail.DirectPath,
		ThumbnailSHA256:     thumbnail.FileSHA256,
		ThumbnailEncSHA256:  thumbnail.FileEncSHA256,
		JPEGThumbnail:       thumbnailByte,

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadVideoMessageFromPath upload video from given path
func (context *Ctx) UploadVideoMessageFromPath(path, caption string) (*waE2E.VideoMessage, error) {
	videoBuff, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer videoBuff.Close()

	videoInfo, _ := videoBuff.Stat()
	videoSize := videoInfo.Size()
	videoBytes := make([]byte, videoSize)

	videoBuffer := bufio.NewReader(videoBuff)
	_, err = videoBuffer.Read(videoBytes)
	if err != nil {
		return nil, err
	}

	mimetypeString := mimetype.Detect(videoBytes)

	thumbnailByte := thumbnail.CreateVideoThumbnail(videoBytes)

	resp, err := context.Client().Upload(context2.Background(), videoBytes, whatsmeow.MediaVideo)
	if err != nil {
		return nil, err
	}

	thumbnail, err := context.Client().Upload(context2.Background(), thumbnailByte, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		thumbnailByte = nil
	}()

	return &waE2E.VideoMessage{
		Caption:  proto.String(caption),
		Mimetype: proto.String(mimetypeString.String()),

		ThumbnailDirectPath: &thumbnail.DirectPath,
		ThumbnailSHA256:     thumbnail.FileSHA256,
		ThumbnailEncSHA256:  thumbnail.FileEncSHA256,
		JPEGThumbnail:       thumbnailByte,

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadVideoMessageFromBytes upload video from given bytes
func (context *Ctx) UploadVideoMessageFromBytes(bytes []byte, caption string) (*waE2E.VideoMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	thumbnailByte := thumbnail.CreateVideoThumbnail(bytes)

	resp, err := context.Client().Upload(context2.Background(), bytes, whatsmeow.MediaVideo)
	if err != nil {
		return nil, err
	}

	thumbnail, err := context.Client().Upload(context2.Background(), thumbnailByte, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		thumbnailByte = nil
		bytes = nil
	}()

	return &waE2E.VideoMessage{
		ContextInfo: util.WithReply(context.MessageEvent()),

		Caption:  proto.String(caption),
		Mimetype: proto.String(mimetypeString.String()),

		ThumbnailDirectPath: &thumbnail.DirectPath,
		ThumbnailSHA256:     thumbnail.FileSHA256,
		ThumbnailEncSHA256:  thumbnail.FileEncSHA256,
		JPEGThumbnail:       thumbnailByte,

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadStickerMessageFromPath upload sticker from given path
func (context *Ctx) UploadStickerMessageFromPath(path string) (*waE2E.StickerMessage, error) {
	stickerBuff, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer stickerBuff.Close()

	stickerInfo, _ := stickerBuff.Stat()
	stickerSize := stickerInfo.Size()
	stickerBytes := make([]byte, stickerSize)

	stickerBuffer := bufio.NewReader(stickerBuff)
	_, err = stickerBuffer.Read(stickerBytes)
	if err != nil {
		return nil, err
	}

	resp, err := context.Client().Upload(context2.Background(), stickerBytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	return &waE2E.StickerMessage{
		Mimetype: proto.String("image/webp"),

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadStickerMessageFromBytes upload sticker from given bytes
func (context *Ctx) UploadStickerMessageFromBytes(bytes []byte) (*waE2E.StickerMessage, error) {
	resp, err := context.Client().Upload(context2.Background(), bytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		bytes = nil
	}()

	return &waE2E.StickerMessage{
		ContextInfo: util.WithReply(context.MessageEvent()),

		Mimetype: proto.String("image/webp"),

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadDocumentMessageFromPath upload document from given path
func (context *Ctx) UploadDocumentMessageFromPath(path, title string) (*waE2E.DocumentMessage, error) {
	documentBuff, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer documentBuff.Close()

	documentInfo, _ := documentBuff.Stat()
	documentSize := documentInfo.Size()
	documentBytes := make([]byte, documentSize)

	documentBuffer := bufio.NewReader(documentBuff)
	_, err = documentBuffer.Read(documentBytes)
	if err != nil {
		return nil, err
	}

	mimetypeString := mimetype.Detect(documentBytes)

	resp, err := context.Client().Upload(context2.Background(), documentBytes, whatsmeow.MediaDocument)
	if err != nil {
		return nil, err
	}

	return &waE2E.DocumentMessage{
		Title:    proto.String(title),
		FileName: proto.String(documentInfo.Name()),
		Mimetype: proto.String(mimetypeString.String()),

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadDocumentMessageFromBytes upload document from given bytes
func (context *Ctx) UploadDocumentMessageFromBytes(bytes []byte, title, filename string) (*waE2E.DocumentMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	resp, err := context.Client().Upload(context2.Background(), bytes, whatsmeow.MediaDocument)
	if err != nil {
		return nil, err
	}

	defer func() {
		bytes = nil
	}()

	return &waE2E.DocumentMessage{
		Title:    proto.String(title),
		FileName: proto.String(filename),
		Mimetype: proto.String(mimetypeString.String()),

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadAudioMessageFromPath upload audio message from given path
func (context *Ctx) UploadAudioMessageFromPath(path string) (*waE2E.AudioMessage, error) {
	audioBuff, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer audioBuff.Close()

	audioInfo, _ := audioBuff.Stat()
	audioSize := audioInfo.Size()
	audioBytes := make([]byte, audioSize)

	audioBuffer := bufio.NewReader(audioBuff)
	_, err = audioBuffer.Read(audioBytes)
	if err != nil {
		return nil, err
	}

	mimetypeString := mimetype.Detect(audioBytes)

	resp, err := context.Client().Upload(context2.Background(), audioBytes, whatsmeow.MediaAudio)
	if err != nil {
		return nil, err
	}

	return &waE2E.AudioMessage{
		Mimetype: proto.String(mimetypeString.String()),

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadAudioMessageFromBytes upload audio from given bytes
func (context *Ctx) UploadAudioMessageFromBytes(bytes []byte) (*waE2E.AudioMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	resp, err := context.Client().Upload(context2.Background(), bytes, whatsmeow.MediaAudio)
	if err != nil {
		return nil, err
	}

	defer func() {
		bytes = nil
	}()

	return &waE2E.AudioMessage{
		Mimetype: proto.String(mimetypeString.String()),

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}
