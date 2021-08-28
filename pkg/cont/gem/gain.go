package gem

import (
	"sort"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleGain(
	gems core.GemStore,
	wallets core.WalletStore,
	system *core.System,
) cont.HandlerFunc {
	return func(r *cont.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		gem, err := From(r, gems)
		if err != nil {
			return err
		}

		if gem.Version >= r.Version {
			return nil
		}

		n := decimal.NewFromInt(int64(len(system.Members)))
		avg := gem.Profit.Div(n).Truncate(8)
		if err := require(avg.IsPositive(), "insufficient-profit"); err != nil {
			return err
		}

		traceID := uuid.Modify(r.TraceID, r.Action.String())
		memo := core.TransferAction{
			ID:     r.FollowID,
			Source: r.Action.String(),
		}.Encode()

		var transfers []*core.Transfer
		for _, member := range copyAndSortMembers(system) {
			t := &core.Transfer{
				TraceID:   uuid.Modify(traceID, member),
				AssetID:   gem.ID,
				Amount:    avg,
				Memo:      memo,
				Threshold: 1,
				Opponents: []string{memo},
			}

			gem.Profit = gem.Profit.Sub(avg)
			transfers = append(transfers, t)
		}

		if err := wallets.CreateTransfers(ctx, transfers); err != nil {
			log.WithError(err).Errorln("wallets.CreateTransfers")
			return err
		}

		if err := gems.Save(ctx, gem, r.Version); err != nil {
			log.WithError(err).Errorln("gems.Save")
			return err
		}

		return nil
	}
}

func copyAndSortMembers(system *core.System) []string {
	members := make([]string, len(system.Members))
	for idx := range system.Members {
		members[idx] = system.Members[idx]
	}

	sort.Slice(members, func(i, j int) bool {
		return members[i] < members[j]
	})

	return members
}
