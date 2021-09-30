package vat

import (
	"time"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/pkg/logger"
)

func HandleLock(
	pools core.PoolStore,
	vaults core.VaultStore,
	assetz core.AssetService,
) cont.HandlerFunc {
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
				(minDur >= 0 && minDur <= dur),
			"bad-data"); err != nil {
			return cont.WithFlag(err, cont.FlagRefund)
		}

		pool, err := pools.Find(ctx, r.AssetID)
		if err != nil {
			log.WithError(err).Errorln("pools.Find")
			return err
		}

		if pool.Version == 0 {
			pool.RewardAt = r.Now
			pool.PardonedAt = r.Now

			asset, err := assetz.Find(ctx, r.AssetID)
			if err != nil {
				log.WithError(err).Errorln("assetz.Find")
				return err
			}

			pool.Price = asset.Price
			pool.Name = asset.Symbol
			pool.Logo = asset.Logo
		}

		vault := &core.Vault{
			CreatedAt:   r.Now,
			UpdatedAt:   r.Now,
			ReleasedAt:  time.Unix(0, 0),
			Version:     r.Version,
			TraceID:     r.TraceID,
			UserID:      r.Sender,
			Status:      core.VaultStatusLocking,
			AssetID:     r.AssetID,
			Duration:    dur,
			MinDuration: minDur,
			Amount:      r.Amount,
			Share:       GetShare(r.Amount, dur),
			LockedPrice: pool.Price,
		}

		if pool.Liquidity.IsPositive() {
			vault.Liquidity = vault.Share.Div(pool.Share.Add(pool.Reward)).Mul(pool.Liquidity).Truncate(12)
		} else {
			vault.Liquidity = vault.Share
		}

		if err := require(vault.Liquidity.IsPositive(), "insufficient-liquidity"); err != nil {
			return cont.WithFlag(err, cont.FlagRefund)
		}

		if err := vaults.Create(ctx, vault); err != nil {
			log.WithError(err).Errorln("vaults.Create")
			return err
		}

		if pool.Version < r.Version {
			pool.Amount = pool.Amount.Add(vault.Amount)
			pool.Share = pool.Share.Add(vault.Share)
			pool.Liquidity = pool.Liquidity.Add(vault.Liquidity)

			if err := pools.Save(ctx, pool, r.Version); err != nil {
				log.WithError(err).Errorln("pools.Save")
				return err
			}
		}

		return nil
	}
}
