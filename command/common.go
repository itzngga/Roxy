package command

import (
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func (runFunc *RunFuncContext) Download(message *waProto.Message, quoted bool) ([]byte, error) {
	var msg *waProto.Message
	if quoted {
		msg = util.ParseQuotedMessage(message)
	} else {
		msg = message
	}

	return runFunc.Client.DownloadAny(msg)
}
