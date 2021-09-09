package views

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/handler/rpc/api"
	"github.com/fox-one/holder/pkg/cont/vat"
)

func Vault(vault *core.Vault, pool *core.Pool) *api.Vault {
	v := &api.Vault{
		Id:         vault.TraceID,
		CreatedAt:  Time(&vault.CreatedAt),
		ReleasedAt: Time(&vault.ReleasedAt),
		// UserId:      vault.UserID,
		Status:        api.Vault_Status(vault.Status),
		AssetId:       vault.AssetID,
		Duration:      vault.Duration,
		MinDuration:   vault.MinDuration,
		Amount:        vault.Amount.String(),
		Share:         vault.Share.String(),
		Liquidity:     vault.Liquidity.String(),
		Reward:        vault.Reward.String(),
		Penalty:       vault.Penalty.String(),
		LockedPrice:   vault.LockedPrice.String(),
		ReleasedPrice: vault.ReleasedPrice.String(),
	}

	if pool != nil {
		if vault.Status == core.VaultStatusLocking {
			v.Reward = vat.GetReward(pool, vault).String()
		}

		v.Pool = Pool(pool)
	}

	return v
}
