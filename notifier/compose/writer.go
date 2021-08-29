package compose

import (
	"encoding/base64"
	"encoding/json"

	"github.com/fox-one/holder/internal/color"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/text/localizer"
)

type Writer struct {
	messageID string
	text      string
	buttons   mixin.AppButtonGroupMessage
	i18n      *localizer.Localizer
}

func (w *Writer) Localize(tmpl string, args ...interface{}) string {
	return w.i18n.LocalizeOr(tmpl, tmpl, args...)
}

func (w *Writer) Text(text string) {
	w.text = text
}

func (w *Writer) Button(label, action string) {
	w.buttons = append(w.buttons, mixin.AppButtonMessage{
		Label:  label,
		Action: action,
		Color:  color.Random(),
	})
}

func (w *Writer) messageRequests() []*mixin.MessageRequest {
	var requests []*mixin.MessageRequest

	if w.text != "" {
		requests = append(requests, &mixin.MessageRequest{
			MessageID: uuid.Modify(w.messageID, "body"),
			Category:  mixin.MessageCategoryPlainText,
			Data:      base64.StdEncoding.EncodeToString([]byte(w.text)),
		})
	}

	if len(w.buttons) > 0 {
		b, _ := json.Marshal(w.buttons)
		requests = append(requests, &mixin.MessageRequest{
			MessageID: uuid.Modify(w.messageID, "buttons"),
			Category:  mixin.MessageCategoryAppButtonGroup,
			Data:      base64.StdEncoding.EncodeToString(b),
		})
	}

	return requests
}
