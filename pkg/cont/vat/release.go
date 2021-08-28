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
	gems core.GemStore,
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

		gem, err := gems.Find(ctx, vault.AssetID)
		if err != nil {
			log.WithError(err).Errorln("gems.Find")
			return err
		}

		if vault.Version < r.Version {
			if err := unlock(r, gem, vault); err != nil {
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

		if gem.Version < r.Version {
			gem.Liquidity = gem.Liquidity.Sub(vault.Liquidity)
			gem.Amount = gem.Amount.Sub(vault.Amount).Sub(vault.Reward)

			if vault.Penalty.IsPositive() {
				gem.Amount = gem.Amount.Add(vault.Penalty)

				v, err := properties.Get(ctx, sys.SystemProfitRateKey)
				if err != nil {
					log.WithError(err).Errorln("properties.Get")
					return err
				}

				rate := number.Decimal(v.String())
				profit := vault.Penalty.Mul(rate).Truncate(8)
				if profit.IsPositive() && profit.LessThanOrEqual(vault.Penalty) {
					gem.Amount = gem.Amount.Sub(profit)
					gem.Profit = gem.Profit.Add(profit)
				}
			}

			if err := gems.Save(ctx, gem, r.Version); err != nil {
				log.WithError(err).Errorln("gems.Save")
				return err
			}
		}

		return nil
	}
}

func unlock(r *cont.Request, gem *core.Gem, vault *core.Vault) error {
	if err := require(vault.Status == core.VaultStatusLocking, "already-unlocked"); err != nil {
		return err
	}

	dur := int64(r.Now.Sub(vault.CreatedAt).Seconds())
	if err := require(dur >= vault.MinDuration, "not-due"); err != nil {
		return err
	}

	if err := require(r.Sender == vault.UserID || dur >= vault.Duration, "not-allowed"); err != nil {
		return err
	}

	vault.Reward = GetReward(gem, vault)
	if remain := vault.Duration - dur; remain > 0 {
		d1 := decimal.NewFromInt(remain)
		d2 := decimal.NewFromInt(vault.Duration)
		p := d1.Div(d2).Mul(vault.Amount).Truncate(8)
		// 扣除本金加奖金
		vault.Penalty = p.Add(vault.Reward)
	}

	vault.Status = core.VaultStatusReleased
	vault.ReleasedAt = r.Now
	return nil
}
