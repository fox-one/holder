package views

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/handler/rpc/api"
)

func Pool(pool *core.Pool) *api.Pool {
	return &api.Pool{
		Id:        pool.ID,
		Name:      pool.Name,
		Logo:      pool.Logo,
		Amount:    pool.Amount.String(),
		Share:     pool.Share.String(),
		Reward:    pool.Reward.String(),
		Liquidity: pool.Liquidity.String(),
		Profit:    pool.Profit.String(),
		Price:     pool.Price.String(),
	}
}
