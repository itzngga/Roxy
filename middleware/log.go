package middleware

import (
	"fmt"
	"github.com/itzngga/goRoxy/internal/handler"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func LogMiddleware(c *whatsmeow.Client, m *events.Message, cmd *handler.Command) bool {
	fmt.Println("\n[CMD] Command : " + cmd.Name)
	return true
}
