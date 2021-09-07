package vat

import (
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
	share := pool.Share
	if vault.Liquidity.LessThan(pool.Liquidity) {
		share = vault.Liquidity.Div(pool.Liquidity).Mul(share).Truncate(12)
	}

	return decimal.Min(
		decimal.Max(
			share.Sub(vault.Share()).Truncate(8),
			decimal.Zero,
		),
		pool.Reward,
	)
}
