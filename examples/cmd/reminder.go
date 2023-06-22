package cmd

import (
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

func init() {
	embed.Commands.Add(reminder)
}

var reminder = &command.Command{
	Name:        "reminder",
	Category:    "UTILITY",
	Description: "Reminder back user",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		var captured string
		command.NewUserQuestion(ctx).
			SetQuestion("Message?", &captured).
			Exec()

		time.Sleep(time.Second * 3)

		//result, err := ctx.Client.DownloadAny(captured)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//res := cli.ExecPipeline("tesseract", result, "stdin", "stdout", "-l", "ind", "--oem", "1", "--psm", "3", "-c", "preserve_interword_spaces=1")

		var bro string
		command.NewUserQuestion(ctx).
			SetReplyQuestion("Bro?", &bro).
			Exec()

		return ctx.GenerateReplyMessage(bro)
	},
}
