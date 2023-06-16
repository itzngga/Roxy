package command

import (
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

// Download message with get quoted message
func (runFunc *RunFuncContext) Download(message *waProto.Message, quoted bool) ([]byte, error) {
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
