package parliament

import (
	"context"
	"fmt"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/mtg"
	"github.com/fox-one/holder/pkg/number"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/shopspring/decimal"
)

func (s *parliament) renderProposalItems(ctx context.Context, action core.Action, data []byte) (items []Item) {
	switch action {
	case core.ActionSysWithdraw:
		var (
			assetID  uuid.UUID
			amount   decimal.Decimal
			opponent uuid.UUID
		)

		_, _ = mtg.Scan(data, &assetID, &amount, &opponent)
		items = []Item{
			{
				Key:    "asset",
				Value:  fmt.Sprintf("%s %s", number.Humanize(amount), s.fetchAssetSymbol(ctx, assetID.String())),
				Action: assetAction(assetID.String()),
			},
			{
				Key:    "opponent",
				Value:  s.fetchUserName(ctx, opponent.String()),
				Action: userAction(opponent.String()),
			},
		}
	case core.ActionSysProperty:
		var key, value string
		_, _ = mtg.Scan(data, &key, &value)
		items = []Item{
			{
				Key:   key,
				Value: value,
			},
		}
	}

	return
}
