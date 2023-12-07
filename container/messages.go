package container

import (
	"context"
	"errors"
	"fmt"
	waProto "github.com/go-whatsapp/whatsmeow/binary/proto"
	waTypes "github.com/go-whatsapp/whatsmeow/types"
	"github.com/go-whatsapp/whatsmeow/types/events"
	"github.com/itzngga/Roxy/util/compress"
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

func (container *Container) InitializeTables() {
	// create chats table
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := container.DB.ExecContext(ctx, CHATS_TABLE)
	if err != nil {
		return
	}
}

func (container *Container) HandleHistorySync(evt *waProto.HistorySync) {
	var currentJID = container.clientJID.String()
	// store status messages
	if len(evt.StatusV3Messages) >= 1 {
		var messages = make([]*events.Message, 0)
		for _, message := range evt.StatusV3Messages {
			// convert the messages
			msg, err := container.ParseWebMessage(waTypes.StatusBroadcastJID, message)
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

		_, err = container.DB.Exec(INSERT_CHATS, currentJID, waTypes.StatusBroadcastJID.String(), result)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
				container.UpsertMessages(waTypes.StatusBroadcastJID, messages)
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
				message, err := container.ParseWebMessage(chatJID, msg.Message)
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

			_, err = container.DB.Exec(INSERT_CHATS, currentJID, chatJID.ToNonAD().String(), result)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
					container.UpsertMessages(chatJID, messages)
				}
				return
			}
		}
	}
}

func (container *Container) HandleMessageUpdates(evt *events.Message) {
	if evt.Message != nil {
		container.UpsertMessages(evt.Info.Chat, []*events.Message{evt})
		return
	}
}

func (container *Container) UpsertMessages(jid waTypes.JID, message []*events.Message) {
	chats := container.GetChatInJID(jid)
	if len(chats) >= 1 {
		chats = append(chats, message...)

		sort.Slice(chats, func(i, j int) bool {
			return chats[i].Info.Timestamp.Before(chats[j].Info.Timestamp)
		})

		result, err := compress.MarshallBrotli(chats)
		if err != nil {
			return
		}

		_, err = container.DB.Exec(UPSERT_CHATS, container.clientJID.String(), jid.ToNonAD().String(), result)
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

		_, err = container.DB.Exec(INSERT_CHATS, container.clientJID.String(), jid.ToNonAD().String(), result)
		if err != nil {
			return
		}
	}
}

func (container *Container) GetAllChats() []*events.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var currentJID = container.clientJID.String()
	rows, err := container.DB.QueryContext(ctx, SELECT_ALL_CHATS, currentJID)
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

func (container *Container) GetChatInJID(jid waTypes.JID) []*events.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var rawMessage []byte
	row := container.DB.QueryRowContext(ctx, SELECT_CHATS_BY_JID, container.clientJID.String(), jid.ToNonAD().String())
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

func (container *Container) GetStatusMessages() []*events.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var rawMessage []byte
	row := container.DB.QueryRowContext(ctx, SELECT_CHATS_BY_JID, container.clientJID.String(), waTypes.StatusBroadcastJID.String())
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

func (container *Container) FindMessageByID(jid waTypes.JID, id string) *events.Message {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var rawMessage []byte
	row := container.DB.QueryRowContext(ctx, SELECT_CHATS_BY_JID, container.clientJID.String(), jid.ToNonAD().String())
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

func (container *Container) ParseWebMessage(chatJID waTypes.JID, webMsg *waProto.WebMessageInfo) (*events.Message, error) {
	var err error
	if chatJID.IsEmpty() {
		chatJID, err = waTypes.ParseJID(webMsg.GetKey().GetRemoteJid())
		if err != nil {
			return nil, fmt.Errorf("no chat JID provided and failed to parse remote JID: %w", err)
		}
	}
	info := waTypes.MessageInfo{
		MessageSource: waTypes.MessageSource{
			Chat:     chatJID,
			IsFromMe: webMsg.GetKey().GetFromMe(),
			IsGroup:  chatJID.Server == waTypes.GroupServer,
		},
		ID:        webMsg.GetKey().GetId(),
		PushName:  webMsg.GetPushName(),
		Timestamp: time.Unix(int64(webMsg.GetMessageTimestamp()), 0),
	}
	if info.IsFromMe {
		info.Sender = container.clientJID.ToNonAD()
		if info.Sender.IsEmpty() {
			return nil, errors.New("the store doesn't contain a device JID")
		}
	} else if chatJID.Server == waTypes.DefaultUserServer || chatJID.Server == waTypes.NewsletterServer {
		info.Sender = chatJID
	} else if webMsg.GetParticipant() != "" {
		info.Sender, err = waTypes.ParseJID(webMsg.GetParticipant())
	} else if webMsg.GetKey().GetParticipant() != "" {
		info.Sender, err = waTypes.ParseJID(webMsg.GetKey().GetParticipant())
	} else {
		return nil, fmt.Errorf("couldn't find sender of message %s", info.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse sender of message %s: %v", info.ID, err)
	}
	evt := &events.Message{
		RawMessage:   webMsg.GetMessage(),
		SourceWebMsg: webMsg,
		Info:         info,
	}
	evt.UnwrapRaw()
	return evt, nil
}
