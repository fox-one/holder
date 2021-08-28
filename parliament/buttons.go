package parliament

import (
	"fmt"

	"github.com/fox-one/holder/internal/color"
	"github.com/fox-one/mixin-sdk-go"
)

func generateButtons(items []Item) mixin.AppButtonGroupMessage {
	var (
		buttons = mixin.AppButtonGroupMessage{}
		idx     int
	)

	for _, item := range items {
		if item.Action == "" {
			continue
		}

		buttons = append(buttons, mixin.AppButtonMessage{
			Label:  item.Value,
			Action: item.Action,
			Color:  color.Random(),
		})

		if idx = idx + 1; idx >= 6 {
			break
		}
	}

	return buttons
}

func assetAction(id string) string {
	return fmt.Sprintf("https://mixin.one/snapshots/%s", id)
}

func userAction(id string) string {
	return mixin.URL.Users(id)
}

func paymentAction(code string) string {
	return mixin.URL.Codes(code)
}
