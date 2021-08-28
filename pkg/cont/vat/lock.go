package vat

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleLock(
	gems core.GemStore,
	vaults core.VaultStore,
	assetz core.AssetService,
) cont.HandlerFunc {
	secondsOfYear := decimal.NewFromInt(365 * 24 * 60 * 60)

	return func(r *cont.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		var (
			dur    int64
			minDur int64
		)

		if err := require(
			r.Scan(&dur, &minDur) == nil &&
				dur > 0 &&
				minDur >= 0 &&
				minDur <= dur,
			"bad-data"); err != nil {
			return cont.WithFlag(err, cont.FlagRefund)
		}

		gem, err := gems.Find(ctx, r.AssetID)
		if err != nil {
			log.WithError(err).Errorln("gems.Find")
			return err
		}

		if gem.Version == 0 {
			asset, err := assetz.Find(ctx, r.AssetID)
			if err != nil {
				log.WithError(err).Errorln("assetz.Find")
				return err
			}

			gem.Name = asset.Name
			gem.Logo = asset.Logo
		}

		vat := &core.Vault{
			CreatedAt:   r.Now,
			UpdatedAt:   r.Now,
			Version:     r.Version,
			TraceID:     r.TraceID,
			UserID:      r.Sender,
			Status:      core.VaultStatusLocking,
			AssetID:     r.AssetID,
			Duration:    dur,
			MinDuration: minDur,
			Amount:      r.Amount,
		}

		if gem.Liquidity.IsPositive() {
			vat.Liquidity = vat.Amount.Div(gem.Amount).Mul(gem.Liquidity).Truncate(8)
		} else {
			vat.Liquidity = decimal.NewFromInt(dur).Div(secondsOfYear).Mul(vat.Amount).Truncate(8)
		}

		if err := require(vat.Liquidity.IsPositive(), "insufficient-liquidity"); err != nil {
			return cont.WithFlag(err, cont.FlagRefund)
		}

		if err := vaults.Create(ctx, vat); err != nil {
			log.WithError(err).Errorln("vaults.Create")
			return err
		}

		if gem.Version < r.Version {
			gem.Amount = gem.Amount.Add(vat.Amount)
			gem.Liquidity = gem.Liquidity.Add(vat.Liquidity)

			if err := gems.Save(ctx, gem, r.Version); err != nil {
				log.WithError(err).Errorln("gems.Save")
				return err
			}
		}

		return nil
	}
}
