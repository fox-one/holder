package compose

import (
	"context"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/text/localizer"
)

func New(
	system *core.System,
	i18n *localizer.Localizer,
	users core.UserStore,
) *Outbox {
	return &Outbox{
		clientID: system.ClientID,
		i18n:     i18n,
		users:    users,
		writers:  map[string]*Writer{},
	}
}

type Outbox struct {
	clientID string
	i18n     *localizer.Localizer
	users    core.UserStore
	writers  map[string]*Writer
}

func (b *Outbox) Write(messageID, receiptID string) *Writer {
	lang := ""
	if user, err := b.users.Find(context.TODO(), receiptID); err == nil {
		lang = user.Lang
	}

	w := &Writer{
		messageID: messageID,
		i18n:      localizer.WithLanguage(b.i18n, lang),
	}

	b.writers[receiptID] = w
	return w
}

func (b *Outbox) Messages() []*core.Message {
	var messages []*core.Message

	for receiptID, writer := range b.writers {
		conversationID := mixin.UniqueConversationID(b.clientID, receiptID)
		for _, req := range writer.messageRequests() {
			req.ConversationID = conversationID
			req.MessageID = uuid.Modify(req.MessageID, receiptID)
			req.RecipientID = receiptID
			messages = append(messages, core.BuildMessage(req))
		}
	}

	return messages
}
