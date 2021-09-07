package pool

import (
	"sort"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleGain(
	pools core.PoolStore,
	wallets core.WalletStore,
	system *core.System,
) cont.HandlerFunc {
	return func(r *cont.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		pool, err := From(r, pools)
		if err != nil {
			return err
		}

		if pool.Version >= r.Version {
			return nil
		}

		n := decimal.NewFromInt(int64(len(system.Members)))
		avg := pool.Profit.Div(n).Truncate(8)
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
				AssetID:   pool.ID,
				Amount:    avg,
				Memo:      memo,
				Threshold: 1,
				Opponents: []string{memo},
			}

			pool.Profit = pool.Profit.Sub(avg)
			transfers = append(transfers, t)
		}

		if err := wallets.CreateTransfers(ctx, transfers); err != nil {
			log.WithError(err).Errorln("wallets.CreateTransfers")
			return err
		}

		if err := pools.Save(ctx, pool, r.Version); err != nil {
			log.WithError(err).Errorln("pools.Save")
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
