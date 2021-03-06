package vat

import (
	"time"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func require(condition bool, msg string) error {
	return cont.Require(condition, "Vat/"+msg)
}

func From(r *cont.Request, vaults core.VaultStore) (*core.Vault, error) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	vat, err := vaults.Find(ctx, id.String())
	if err != nil {
		log.WithError(err).Errorln("vaults.Find")
		return nil, err
	}

	if err := require(vat.ID > 0, "not init"); err != nil {
		return nil, err
	}

	return vat, nil
}

func GetReward(pool *core.Pool, vault *core.Vault) decimal.Decimal {
	if !pool.RewardAt.IsZero() && vault.CreatedAt.Sub(pool.RewardAt) >= 0 {
		return decimal.Zero
	}

	share := pool.Share.Add(pool.Reward)
	if vault.Liquidity.LessThan(pool.Liquidity) {
		share = vault.Liquidity.Div(pool.Liquidity).Mul(share).Truncate(12)
	}

	return decimal.Min(
		decimal.Max(
			share.Sub(vault.Share).Truncate(8),
			decimal.Zero,
		),
		pool.Reward,
	)
}

func GetShare(amount decimal.Decimal, dur int64) decimal.Decimal {
	year := time.Hour * 24 * 365
	secondsOfYear := year.Milliseconds() / 1000
	d := decimal.NewFromInt(dur).Div(decimal.NewFromInt(secondsOfYear))
	return amount.Mul(d).Truncate(12)
}
