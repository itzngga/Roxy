package core

import (
	"github.com/itzngga/goRoxy/command"
	"github.com/itzngga/goRoxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func (m *Muxer) GenerateHelpButton() {
	m.AddCommand(&command.Command{
		Name:    "help",
		Aliases: []string{"menu"},
		Cache:   true,
		RunFunc: func(c *whatsmeow.Client, args command.RunFuncArgs) *waProto.Message {
			id, _ := args.Locals.Load("uid")

			sections := make([]*waProto.ListMessage_Section, 0)
			sections = append(sections, util.CreateSectionList("Umum", util.CreateSectionRow("Daftar Perintah", "/help", util.CreateButtonID(id, "/help"))))

			m.Categories.Range(func(ctKey string, category string) bool {
				rows := make([]*waProto.ListMessage_Row, 0)
				m.Commands.Range(func(cmdKey string, cmd *command.Command) bool {
					if !cmd.HideFromHelp && cmd.Category == category {
						for _, row := range rows {
							if *row.Description == "/"+cmd.Name {
								return true
							}
						}
						rows = append(rows, util.CreateSectionRow(cmd.Description, "/"+cmd.Name, util.CreateButtonID(id, "/"+cmd.Name)))
					}
					return true
				})
				sections = append(sections, util.CreateSectionList(category, rows...))
				return true
			})

			return util.GenerateListMessage(m.Options.HelpTitle, m.Options.HelpDescription, "Lihat", m.Options.HelpFooter, sections...)
		},
	})
}
