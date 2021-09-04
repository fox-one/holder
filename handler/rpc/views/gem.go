package views

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/handler/rpc/api"
)

func Gem(gem *core.Gem) *api.Gem {
	return &api.Gem{
		Id:        gem.ID,
		Name:      gem.Name,
		Logo:      gem.Logo,
		Amount:    gem.Amount.String(),
		Reward:    gem.Reward.String(),
		Liquidity: gem.Liquidity.String(),
		Profit:    gem.Profit.String(),
	}
}
