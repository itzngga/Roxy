package util

import (
	"bufio"
	"context"
	"github.com/gabriel-vasile/mimetype"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
	"os"
)

func UploadImageMessageFromPath(c *whatsmeow.Client, path, caption string) (*waProto.ImageMessage, error) {
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

	resp, err := c.Upload(context.Background(), imageBytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	return &waProto.ImageMessage{
		Caption:  proto.String(caption),
		Mimetype: proto.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

func UploadImageMessageFromBytes(c *whatsmeow.Client, bytes []byte, caption string) (*waProto.ImageMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	resp, err := c.Upload(context.Background(), bytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, err
	}

	return &waProto.ImageMessage{
		Caption:  proto.String(caption),
		Mimetype: proto.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

func UploadVideoMessageFromPath(c *whatsmeow.Client, path, caption string) (*waProto.VideoMessage, error) {
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

	resp, err := c.Upload(context.Background(), videoBytes, whatsmeow.MediaVideo)

	if err != nil {
		return nil, err
	}

	return &waProto.VideoMessage{
		Caption:  proto.String(caption),
		Mimetype: proto.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

func UploadVideoMessageFromBytes(c *whatsmeow.Client, bytes []byte, caption string) (*waProto.VideoMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	resp, err := c.Upload(context.Background(), bytes, whatsmeow.MediaVideo)
	if err != nil {
		return nil, err
	}

	return &waProto.VideoMessage{
		Caption:  proto.String(caption),
		Mimetype: proto.String(mimetypeString.String()),

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}, nil
}

func UploadStickerMessageFromPath(c *whatsmeow.Client, path string) (*waProto.StickerMessage, error) {
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

	resp, err := c.Upload(context.Background(), stickerBytes, whatsmeow.MediaImage)

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

func UploadStickerMessageFromBytes(c *whatsmeow.Client, bytes []byte) (*waProto.StickerMessage, error) {
	resp, err := c.Upload(context.Background(), bytes, whatsmeow.MediaImage)
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

func UploadDocumentMessageFromPath(c *whatsmeow.Client, path, title string) (*waProto.DocumentMessage, error) {
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

	resp, err := c.Upload(context.Background(), documentBytes, whatsmeow.MediaDocument)
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

func UploadDocumentMessageFromBytes(c *whatsmeow.Client, bytes []byte, title, filename string) (*waProto.DocumentMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	resp, err := c.Upload(context.Background(), bytes, whatsmeow.MediaDocument)

	if err != nil {
		return nil, err
	}

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

func UploadAudioMessageFromPath(c *whatsmeow.Client, path string) (*waProto.AudioMessage, error) {
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

	resp, err := c.Upload(context.Background(), audioBytes, whatsmeow.MediaAudio)

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

func UploadAudioMessageFromBytes(c *whatsmeow.Client, bytes []byte) (*waProto.AudioMessage, error) {
	mimetypeString := mimetype.Detect(bytes)

	resp, err := c.Upload(context.Background(), bytes, whatsmeow.MediaAudio)
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
