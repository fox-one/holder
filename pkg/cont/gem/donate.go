package gem

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/cont/sys"
	"github.com/fox-one/holder/pkg/mtg/types"
	"github.com/fox-one/holder/pkg/number"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
)

func HandleDonate(gems core.GemStore, properties property.Store) cont.HandlerFunc {
	return func(r *cont.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		gem, err := From(r.WithBody(types.UUID(r.AssetID)), gems)
		if err != nil {
			return err
		}

		if err := require(gem.Liquidity.IsPositive(), "empty"); err != nil {
			return cont.WithFlag(err, cont.FlagRefund)
		}

		if gem.Version < r.Version {
			gem.Amount = gem.Amount.Add(r.Amount)

			v, err := properties.Get(ctx, sys.SystemDonateFeeRate)
			if err != nil {
				log.WithError(err).Errorln("properties.Get")
				return err
			}

			rate := number.Decimal(v.String())
			if fee := r.Amount.Mul(rate).Truncate(8); fee.IsPositive() && fee.LessThanOrEqual(r.Amount) {
				gem.Amount = gem.Amount.Sub(fee)
				gem.Profit = gem.Profit.Add(fee)
			}

			if err := gems.Save(ctx, gem, r.Version); err != nil {
				log.WithError(err).Errorln("gems.Save")
				return err
			}
		}

		return nil
	}
}
