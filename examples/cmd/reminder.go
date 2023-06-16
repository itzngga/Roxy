package cmd

import (
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/util/cli"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"log"
)

func init() {
	embed.Commands.Add(reminder)
}

var reminder = &command.Command{
	Name:        "reminder",
	Category:    "UTILITY",
	Description: "Reminder back user",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		var captured *waProto.Message
		command.NewUserQuestion(ctx).
			CaptureMediaQuestion("Message?", &captured).
			Exec()

		result, err := ctx.Client.DownloadAny(captured)
		if err != nil {
			log.Fatal(err)
		}
		res := cli.ExecPipeline("tesseract", result, "stdin", "stdout", "-l", "ind", "--oem", "1", "--psm", "3", "-c", "preserve_interword_spaces=1")

		return ctx.GenerateReplyMessage(string(res))
	},
}
