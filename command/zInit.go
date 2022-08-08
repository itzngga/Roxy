package command

import "github.com/itzngga/goRoxy/internal/handler"

var Commands []*handler.Command

func GenerateAllCommands() {
	HiCommand()
	StickerCommand()
	ButtonCommand()
	HydratedCommand()
	ImageButtonCommand()
	UrlImageButtonCommand()
}

func AddCommand(command *handler.Command) {
	Commands = append(Commands, command)
}
