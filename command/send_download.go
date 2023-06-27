package command

import (
	"bytes"
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"io"
	"os"
)

// Download message with get quoted message
func (runFunc *RunFuncContext) Download(quoted bool) ([]byte, error) {
	var msg *waProto.Message
	if quoted {
		result := util.ParseQuotedMessage(runFunc.Message)
		if result != nil {
			msg = result
		} else {
			msg = runFunc.Message
		}
	} else {
		msg = runFunc.Message
	}

	return runFunc.Client.DownloadAny(msg)
}

// DownloadToFile download message to file with quoted message
func (runFunc *RunFuncContext) DownloadToFile(quoted bool, fileName string) (*os.File, error) {
	var msg *waProto.Message
	if quoted {
		result := util.ParseQuotedMessage(runFunc.Message)
		if result != nil {
			msg = result
		} else {
			msg = runFunc.Message
		}
	} else {
		msg = runFunc.Message
	}

	data, err := runFunc.Client.DownloadAny(msg)
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
func (runFunc *RunFuncContext) DownloadMessage(message *waProto.Message, quoted bool) ([]byte, error) {
	var msg *waProto.Message
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

	return runFunc.Client.DownloadAny(msg)
}

// DownloadMessageToFile download with given message to file
func (runFunc *RunFuncContext) DownloadMessageToFile(message *waProto.Message, quoted bool, fileName string) (*os.File, error) {
	var msg *waProto.Message
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

	data, err := runFunc.Client.DownloadAny(msg)
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
func (runFunc *RunFuncContext) GetDownloadable(quoted bool) *waProto.Message {
	var msg *waProto.Message
	if quoted {
		result := util.ParseQuotedMessage(runFunc.Message)
		if result != nil {
			msg = result
		} else {
			msg = runFunc.Message
		}
	} else {
		msg = runFunc.Message
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
func (runFunc *RunFuncContext) GetDownloadableMessage(message *waProto.Message, quoted bool) *waProto.Message {
	var msg *waProto.Message
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
