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

func GetReward(gem *core.Gem, vat *core.Vault) decimal.Decimal {
	amount := gem.Amount
	if vat.Liquidity.LessThan(gem.Liquidity) {
		amount = vat.Liquidity.Div(gem.Liquidity).Mul(amount).Truncate(8)
	}

	return decimal.Min(
		decimal.Max(
			amount.Sub(vat.Amount),
			decimal.Zero,
		),
		gem.Reward,
	)
}
