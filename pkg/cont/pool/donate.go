package pool

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/cont/sys"
	"github.com/fox-one/holder/pkg/mtg/types"
	"github.com/fox-one/holder/pkg/number"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
)

func HandleDonate(pools core.PoolStore, properties property.Store) cont.HandlerFunc {
	return func(r *cont.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		pool, err := From(r.WithBody(types.UUID(r.AssetID)), pools)
		if err != nil {
			return err
		}

		if err := require(pool.Liquidity.IsPositive(), "empty"); err != nil {
			return cont.WithFlag(err, cont.FlagRefund)
		}

		if pool.Version < r.Version {
			pool.Amount = pool.Amount.Add(r.Amount)
			pool.Reward = pool.Reward.Add(r.Amount)
			pool.Share = pool.Share.Add(r.Amount)

			v, err := properties.Get(ctx, sys.SystemDonateFeeRate)
			if err != nil {
				log.WithError(err).Errorln("properties.Get")
				return err
			}

			rate := number.Decimal(v.String())
			if fee := r.Amount.Mul(rate).Truncate(8); fee.IsPositive() && fee.LessThanOrEqual(r.Amount) {
				pool.Amount = pool.Amount.Sub(fee)
				pool.Reward = pool.Reward.Sub(fee)
				pool.Share = pool.Share.Sub(fee)
				pool.Profit = pool.Profit.Add(fee)
			}

			if err := pools.Save(ctx, pool, r.Version); err != nil {
				log.WithError(err).Errorln("pools.Save")
				return err
			}
		}

		return nil
	}
}
