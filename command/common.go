package command

import (
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
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
