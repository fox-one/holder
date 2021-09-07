package vat

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/pkg/logger"
)

func HandleLock(
	pools core.PoolStore,
	vaults core.VaultStore,
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

		vault := &core.Vault{
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

		share := vault.Share()
		if pool.Liquidity.IsPositive() {
			vault.Liquidity = share.Div(pool.Share).Mul(pool.Liquidity).Truncate(12)
		} else {
			vault.Liquidity = share
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
			pool.Share = pool.Share.Add(share)
			pool.Liquidity = pool.Liquidity.Add(vault.Liquidity)

			if err := pools.Save(ctx, pool, r.Version); err != nil {
				log.WithError(err).Errorln("pools.Save")
				return err
			}
		}

		return nil
	}
}
