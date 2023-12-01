package context

import (
	"bufio"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"github.com/go-whatsapp/whatsmeow"
	waProto "github.com/go-whatsapp/whatsmeow/binary/proto"
	"github.com/itzngga/Roxy/util"
	"github.com/itzngga/Roxy/util/thumbnail"
	"google.golang.org/protobuf/proto"
	"os"

	context2 "context"
)

var imageExtensions = []string{".png", ".jpg", ".webp", ".gif", ".bmp", ".ico", ".svg"}
var videoExtensions = []string{".mp4", ".mov", ".mpeg", ".webp", ".3gp", ".avi", ".mkv"}
var audioExtensions = []string{".mp3", ".ogg", ".m4a", ".wav", ".flac"}

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

	for _, audioExtension := range audioExtensions {
		if extension == audioExtension {
			return context.UploadAudioMessageFromBytes(bytes)
		}
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
func (context *Ctx) UploadImageMessageFromPath(path, caption string) (*waProto.ImageMessage, error) {
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

	return &waProto.ImageMessage{
		Caption: proto.String(caption),

		Mimetype: proto.String(mimetypeString.String()),

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
func (context *Ctx) UploadImageMessageFromBytes(bytes []byte, caption string) (*waProto.ImageMessage, error) {
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

	return &waProto.ImageMessage{
		ContextInfo: util.WithReply(context.MessageEvent()),

		Caption:  proto.String(caption),
		Mimetype: proto.String(mimetypeString.String()),

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
func (context *Ctx) UploadVideoMessageFromPath(path, caption string) (*waProto.VideoMessage, error) {
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

	return &waProto.VideoMessage{
		Caption:  proto.String(caption),
		Mimetype: proto.String(mimetypeString.String()),

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
func (context *Ctx) UploadVideoMessageFromBytes(bytes []byte, caption string) (*waProto.VideoMessage, error) {
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

	return &waProto.VideoMessage{
		ContextInfo: util.WithReply(context.MessageEvent()),

		Caption:  proto.String(caption),
		Mimetype: proto.String(mimetypeString.String()),

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
func (context *Ctx) UploadStickerMessageFromPath(path string) (*waProto.StickerMessage, error) {
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

	resp, err := context.Client().Upload(context2.Background(), stickerBytes, whatsmeow.MediaImage)

	if err != nil {
		return nil, err
	}

	return &waProto.StickerMessage{
		Mimetype: proto.String("image/webp"),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadStickerMessageFromBytes upload sticker from given bytes
func (context *Ctx) UploadStickerMessageFromBytes(bytes []byte) (*waProto.StickerMessage, error) {
	resp, err := context.Client().Upload(context2.Background(), bytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	defer func() {
		bytes = nil
	}()

	return &waProto.StickerMessage{
		ContextInfo: util.WithReply(context.MessageEvent()),

		Mimetype: proto.String("image/webp"),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadDocumentMessageFromPath upload document from given path
func (context *Ctx) UploadDocumentMessageFromPath(path, title string) (*waProto.DocumentMessage, error) {
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

	resp, err := context.Client().Upload(context2.Background(), documentBytes, whatsmeow.MediaDocument)
	if err != nil {
		return nil, err
	}

	return &waProto.DocumentMessage{
		Title:    proto.String(title),
		FileName: proto.String(documentInfo.Name()),
		Mimetype: proto.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadDocumentMessageFromBytes upload document from given bytes
func (context *Ctx) UploadDocumentMessageFromBytes(bytes []byte, title, filename string) (*waProto.DocumentMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	resp, err := context.Client().Upload(context2.Background(), bytes, whatsmeow.MediaDocument)

	if err != nil {
		return nil, err
	}

	defer func() {
		bytes = nil
	}()

	return &waProto.DocumentMessage{
		Title:    proto.String(title),
		FileName: proto.String(filename),
		Mimetype: proto.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadAudioMessageFromPath upload audio message from given path
func (context *Ctx) UploadAudioMessageFromPath(path string) (*waProto.AudioMessage, error) {
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

	resp, err := context.Client().Upload(context2.Background(), audioBytes, whatsmeow.MediaAudio)

	if err != nil {
		return nil, err
	}

	return &waProto.AudioMessage{
		Mimetype: proto.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

// UploadAudioMessageFromBytes upload audio from given bytes
func (context *Ctx) UploadAudioMessageFromBytes(bytes []byte) (*waProto.AudioMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	resp, err := context.Client().Upload(context2.Background(), bytes, whatsmeow.MediaAudio)
	if err != nil {
		return nil, err
	}

	defer func() {
		bytes = nil
	}()

	return &waProto.AudioMessage{
		Mimetype: proto.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}
