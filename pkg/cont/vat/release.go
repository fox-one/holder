package vat

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/cont/sys"
	"github.com/fox-one/holder/pkg/number"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
	"github.com/shopspring/decimal"
)

func HandleRelease(
	pools core.PoolStore,
	vaults core.VaultStore,
	wallets core.WalletStore,
	properties property.Store,
) cont.HandlerFunc {
	return func(r *cont.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		vault, err := From(r, vaults)
		if err != nil {
			return err
		}

		pool, err := pools.Find(ctx, vault.AssetID)
		if err != nil {
			log.WithError(err).Errorln("pools.Find")
			return err
		}

		pool.Reform(vault)

		if vault.Version < r.Version {
			if err := unlock(r, pool, vault); err != nil {
				return err
			}

			memo := core.TransferAction{
				ID:     r.FollowID,
				Source: r.Action.String(),
			}.Encode()

			t := &core.Transfer{
				TraceID:   uuid.Modify(r.TraceID, r.Action.String()),
				AssetID:   vault.AssetID,
				Amount:    vault.Amount.Add(vault.Reward).Sub(vault.Penalty),
				Memo:      memo,
				Threshold: 1,
				Opponents: []string{vault.UserID},
			}

			if err := wallets.CreateTransfers(ctx, []*core.Transfer{t}); err != nil {
				log.WithError(err).Errorln("wallets.CreateTransfers")
				return err
			}

			if err := vaults.Update(ctx, vault, r.Version); err != nil {
				log.WithError(err).Errorln("vaults.Update")
				return err
			}
		}

		if pool.Version < r.Version {
			pool.Liquidity = pool.Liquidity.Sub(vault.Liquidity)
			pool.Share = pool.Share.Sub(vault.Share)
			pool.Amount = pool.Amount.Sub(vault.Amount).Sub(vault.Reward)
			pool.Reward = pool.Reward.Sub(vault.Reward)

			if vault.Penalty.IsPositive() {
				pool.Amount = pool.Amount.Add(vault.Penalty)
				pool.Reward = pool.Reward.Add(vault.Penalty)
				pool.RewardAt = r.Now

				v, err := properties.Get(ctx, sys.SystemProfitRateKey)
				if err != nil {
					log.WithError(err).Errorln("properties.Get")
					return err
				}

				rate := number.Decimal(v.String())
				profit := vault.Penalty.Mul(rate).Truncate(8)
				if profit.IsPositive() && profit.LessThanOrEqual(vault.Penalty) {
					pool.Amount = pool.Amount.Sub(profit)
					pool.Reward = pool.Reward.Sub(profit)
					pool.Profit = pool.Profit.Add(profit)
				}
			}

			if err := pools.Save(ctx, pool, r.Version); err != nil {
				log.WithError(err).Errorln("pools.Save")
				return err
			}
		}

		return nil
	}
}

func unlock(r *cont.Request, pool *core.Pool, vault *core.Vault) error {
	if err := require(r.Sender == vault.UserID, "not-allowed"); err != nil {
		return err
	}

	if err := require(vault.Status == core.VaultStatusLocking, "already-unlocked"); err != nil {
		return err
	}

	dur := int64(r.Now.Sub(vault.CreatedAt).Seconds())
	if err := require(dur >= vault.MinDuration, "not-due"); err != nil {
		return err
	}

	vault.Reward = GetReward(pool, vault)
	if remain := vault.Duration - dur; remain > 0 {
		d1 := decimal.NewFromInt(remain)
		d2 := decimal.NewFromInt(vault.Duration)
		p := d1.Div(d2).Mul(vault.Amount).Truncate(8)
		// 扣除本金加奖金
		vault.Penalty = p.Add(vault.Reward)
	}

	vault.Status = core.VaultStatusReleased
	vault.ReleasedAt = r.Now
	vault.ReleasedPrice = pool.Price
	return nil
}
