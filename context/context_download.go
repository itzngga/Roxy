package context

import (
	"bytes"
	"io"
	"os"

	"github.com/itzngga/Roxy/util"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

// Download message with get quoted message
func (context *Ctx) Download(quoted bool) ([]byte, error) {
	var msg *waE2E.Message
	if quoted {
		result := util.ParseQuotedMessage(context.Message())
		if result != nil {
			msg = result
		} else {
			msg = context.Message()
		}
	} else {
		msg = context.Message()
	}

	return context.Client().DownloadAny(msg)
}

// DownloadToFile download message to file with quoted message
func (context *Ctx) DownloadToFile(quoted bool, fileName string) (*os.File, error) {
	var msg *waE2E.Message
	if quoted {
		result := util.ParseQuotedMessage(context.Message())
		if result != nil {
			msg = result
		} else {
			msg = context.Message()
		}
	} else {
		msg = context.Message()
	}

	data, err := context.Client().DownloadAny(msg)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)
	io.Copy(file, reader)

	defer func() {
		data = nil
		reader = nil
	}()

	return file, nil
}

// DownloadMessage download with given message
func (context *Ctx) DownloadMessage(message *waE2E.Message, quoted bool) ([]byte, error) {
	var msg *waE2E.Message
	if quoted {
		result := util.ParseQuotedMessage(message)
		if result != nil {
			msg = result
		} else {
			msg = message
		}
	} else {
		msg = message
	}

	return context.Client().DownloadAny(msg)
}

// DownloadMessageToFile download with given message to file
func (context *Ctx) DownloadMessageToFile(message *waE2E.Message, quoted bool, fileName string) (*os.File, error) {
	var msg *waE2E.Message
	if quoted {
		result := util.ParseQuotedMessage(message)
		if result != nil {
			msg = result
		} else {
			msg = message
		}
	} else {
		msg = message
	}

	data, err := context.Client().DownloadAny(msg)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)
	io.Copy(file, reader)

	defer func() {
		data = nil
		reader = nil
	}()

	return file, nil
}

// GetDownloadable get downloadable type
func (context *Ctx) GetDownloadable(quoted bool) *waE2E.Message {
	var msg *waE2E.Message
	if quoted {
		result := util.ParseQuotedMessage(context.Message())
		if result != nil {
			msg = result
		} else {
			msg = context.Message()
		}
	} else {
		msg = context.Message()
	}

	switch {
	case msg.ImageMessage != nil:
		return msg
	case msg.VideoMessage != nil:
		return msg
	case msg.AudioMessage != nil:
		return msg
	case msg.DocumentMessage != nil:
		return msg
	case msg.StickerMessage != nil:
		return msg
	default:
		return nil
	}
}

// GetDownloadableMessage get downloadable with given message
func (context *Ctx) GetDownloadableMessage(message *waE2E.Message, quoted bool) *waE2E.Message {
	var msg *waE2E.Message
	if quoted {
		result := util.ParseQuotedMessage(message)
		if result != nil {
			msg = result
		} else {
			msg = message
		}
	} else {
		msg = message
	}

	switch {
	case msg.ImageMessage != nil:
		return msg
	case msg.VideoMessage != nil:
		return msg
	case msg.AudioMessage != nil:
		return msg
	case msg.DocumentMessage != nil:
		return msg
	case msg.StickerMessage != nil:
		return msg
	default:
		return nil
	}
}
