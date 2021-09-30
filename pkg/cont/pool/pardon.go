package pool

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/mtg/types"
	"github.com/fox-one/pkg/logger"
)

func HandlePardon(pools core.PoolStore) cont.HandlerFunc {
	return func(r *cont.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		pool, err := From(r.WithBody(types.UUID(r.AssetID)), pools)
		if err != nil {
			return err
		}

		if pool.Version >= r.Version {
			return nil
		}

		pool.PardonedAt = r.Now
		if err := pools.Save(ctx, pool, r.Version); err != nil {
			log.WithError(err).Errorln("pools.Save")
			return err
		}

		return nil
	}
}
