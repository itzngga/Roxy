package command

import (
	"fmt"
	"github.com/itzngga/goRoxy/internal/handler"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"time"
)

func HiCommand(c *whatsmeow.Client, m *events.Message, cmd *handler.Command) *waProto.Message {
	t := time.Now()
	util.SendReplyMessage(c, m, "testing a...")
	return util.SendReplyText(m, fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
}
