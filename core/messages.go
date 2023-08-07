package core

import (
	"context"
	"github.com/itzngga/Roxy/util/compress"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"sort"
	"strings"
	"time"
)

const CHATS_TABLE = `CREATE TABLE IF NOT EXISTS whatsmeow_chats
(
    our_jid       TEXT,
	chat_jid      TEXT,
    messages      BLOB,

   PRIMARY KEY (our_jid, chat_jid),
   FOREIGN KEY (our_jid) REFERENCES whatsmeow_device(jid) ON DELETE CASCADE ON UPDATE CASCADE
);`
const INSERT_CHATS = `INSERT INTO whatsmeow_chats (our_jid, chat_jid, messages) VALUES ($1, $2, $3)`
const UPSERT_CHATS = `INSERT INTO whatsmeow_chats (our_jid, chat_jid, messages) VALUES ($1, $2, $3) ON CONFLICT (our_jid, chat_jid) DO UPDATE SET messages = excluded.messages`
const SELECT_CHATS_BY_JID = `SELECT messages FROM whatsmeow_chats WHERE our_jid = $1 AND chat_jid = $2`
const SELECT_ALL_CHATS = `SELECT messages FROM whatsmeow_chats WHERE our_jid = $1`
const UPDATE_MESSAGES_IN_CHAT = `UPDATE whatsmeow_chats SET messages = $1 WHERE our_jid = $1 AND chat_jid = $2`
const DELETE_ALL_CHATS = `DELETE FROM whatsmeow_chats WHERE our_jid = $1`

func (app *App) handleHistorySync(evt *waProto.HistorySync) {
	var currentJID = app.client.Store.ID.String()
	// store status messages
	if len(evt.StatusV3Messages) >= 1 {
		var messages = make([]*events.Message, 0)
		for _, message := range evt.StatusV3Messages {
			// convert the messages
			msg, err := app.client.ParseWebMessage(waTypes.StatusBroadcastJID, message)
			if err != nil {
				continue
			}
			messages = append(messages, msg)
		}

		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Info.Timestamp.Before(messages[j].Info.Timestamp)
		})

		result, err := compress.MarshallBrotli(messages)
		if err != nil {
			return
		}

		_, err = app.sqlDB.Exec(INSERT_CHATS, currentJID, waTypes.StatusBroadcastJID.String(), result)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
				app.upsertMessages(waTypes.StatusBroadcastJID, messages)
			}
			return
		}
	}

	if len(evt.Conversations) >= 1 {
		for _, conversation := range evt.Conversations {
			chatJID, _ := waTypes.ParseJID(conversation.GetId())
			if chatJID.IsEmpty() || chatJID == waTypes.StatusBroadcastJID {
				continue
			}

			if len(conversation.GetMessages()) <= 0 {
				return
			}

			// convert the messages
			var messages = make([]*events.Message, 0)
			for _, msg := range conversation.GetMessages() {
				message, err := app.client.ParseWebMessage(chatJID, msg.Message)
				if err != nil {
					continue
				}
				messages = append(messages, message)
			}

			sort.Slice(messages, func(i, j int) bool {
				return messages[i].Info.Timestamp.Before(messages[j].Info.Timestamp)
			})

			result, err := compress.MarshallBrotli(messages)
			if err != nil {
				continue
			}

			_, err = app.sqlDB.Exec(INSERT_CHATS, currentJID, chatJID.ToNonAD().String(), result)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
					app.upsertMessages(chatJID, messages)
				}
				return
			}
		}
	}
}

func (app *App) handleMessageUpdates(evt *events.Message) {
	if evt.Message != nil {
		app.upsertMessages(evt.Info.Chat, []*events.Message{evt})
		return
	}
}

func (app *App) upsertMessages(jid waTypes.JID, message []*events.Message) {
	chats := app.getChatInJID(jid)
	if len(chats) >= 1 {
		chats = append(chats, message...)

		sort.Slice(chats, func(i, j int) bool {
			return chats[i].Info.Timestamp.Before(chats[j].Info.Timestamp)
		})

		result, err := compress.MarshallBrotli(chats)
		if err != nil {
			return
		}

		_, err = app.sqlDB.Exec(UPSERT_CHATS, app.client.Store.ID.String(), jid.ToNonAD().String(), result)
		if err != nil {
			return
		}
	} else {
		sort.Slice(message, func(i, j int) bool {
			return message[i].Info.Timestamp.Before(message[j].Info.Timestamp)
		})

		result, err := compress.MarshallBrotli(message)
		if err != nil {
			return
		}

		_, err = app.sqlDB.Exec(INSERT_CHATS, app.client.Store.ID.String(), jid.ToNonAD().String(), result)
		if err != nil {
			return
		}
	}
}

func (app *App) getAllChats() []*events.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var currentJID = app.client.Store.ID.String()
	rows, err := app.sqlDB.QueryContext(ctx, SELECT_ALL_CHATS, currentJID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var model = make([]*events.Message, 0)
	for rows.Next() {
		var rawMessage []byte
		err = rows.Scan(&rawMessage)
		if err != nil {
			continue
		}

		var message = make([]*events.Message, 0)
		err = compress.UnmarshallBrotli(rawMessage, &message)
		if err != nil {
			continue
		}

		model = append(model, message...)
		rawMessage = nil
	}

	if err := rows.Close(); err != nil {
		return nil
	}
	if err := rows.Err(); err != nil {
		return nil
	}

	return model
}

func (app *App) getChatInJID(jid waTypes.JID) []*events.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var rawMessage []byte
	row := app.sqlDB.QueryRowContext(ctx, SELECT_CHATS_BY_JID, app.client.Store.ID.String(), jid.ToNonAD().String())
	err := row.Scan(&rawMessage)
	if err != nil {
		return nil
	}

	defer func() {
		rawMessage = nil
	}()

	var message = make([]*events.Message, 0)
	err = compress.UnmarshallBrotli(rawMessage, &message)
	if err != nil {
		return nil
	}

	return message
}

func (app *App) getStatusMessages() []*events.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var rawMessage []byte
	row := app.sqlDB.QueryRowContext(ctx, SELECT_CHATS_BY_JID, app.client.Store.ID.String(), waTypes.StatusBroadcastJID.String())
	err := row.Scan(&rawMessage)
	if err != nil {
		return nil
	}

	defer func() {
		rawMessage = nil
	}()

	var message = make([]*events.Message, 0)
	err = compress.UnmarshallBrotli(rawMessage, &message)
	if err != nil {
		return nil
	}

	return message
}

func (app *App) findMessageByID(jid waTypes.JID, id string) *events.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var rawMessage []byte
	row := app.sqlDB.QueryRowContext(ctx, SELECT_CHATS_BY_JID, app.client.Store.ID.String(), jid.ToNonAD().String())
	err := row.Scan(&rawMessage)
	if err != nil {
		return nil
	}

	var message = make([]*events.Message, 0)
	err = compress.UnmarshallBrotli(rawMessage, &message)
	if err != nil {
		return nil
	}

	defer func() {
		rawMessage = nil
		message = nil
	}()

	for _, message := range message {
		if message.Info.ID == id {
			return message
		}
	}

	return nil
}

func (app *App) initializeTables() {
	// create chats table
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := app.sqlDB.ExecContext(ctx, CHATS_TABLE)
	if err != nil {
		return
	}
}
