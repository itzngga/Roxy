package command

import (
	"fmt"
	"github.com/itzngga/goRoxy/internal/handler"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"time"
)

func TestSpeedCommand(c *whatsmeow.Client, m *events.Message) *waProto.Message {
	t := time.Now()
	handler.SendReplyMessage(c, m, "testing a...")
	return handler.SendReplyText(m, fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
}
