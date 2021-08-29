package notifier

import (
	"github.com/fox-one/holder/internal/color"
	"github.com/fox-one/mixin-sdk-go"
)

type TxData struct {
	Action  string
	Message string
	Lines   []string
	Buttons mixin.AppButtonGroupMessage
}

func (data *TxData) AddLine(line string) {
	data.Lines = append(data.Lines, line)
}

func (data *TxData) AddButton(label, action string) {
	data.Buttons = append(data.Buttons, mixin.AppButtonMessage{
		Label:  label,
		Action: action,
		Color:  color.Random(),
	})
}
