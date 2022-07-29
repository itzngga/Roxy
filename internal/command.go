package internal

import (
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/internal/handler"
)

func (b Base) InitializeCommands() {
	handler.AddCommand(&handler.Command{
		Name:        "tes",
		Aliases:     []string{"hi"},
		Description: "A Fucking Test",
		Category:    handler.MiscCategory,
		RunFunc:     command.TestSpeedCommand,
	})
	handler.AddCommand(
		&handler.Command{
			Name:        "sticker",
			Aliases:     []string{"stkr", "stiker"},
			Category:    handler.UtilitiesCategory,
			Description: "Create sticker from image or video",
			RunFunc:     command.StickerCommand,
		})
}
