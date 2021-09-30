package pool

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/pkg/logger"
)

func HandlePardonAll(pools core.PoolStore) cont.HandlerFunc {
	return func(r *cont.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		all, err := pools.List(ctx)
		if err != nil {
			log.WithError(err).Errorln("pools.List")
			return err
		}

		for _, pool := range all {
			if pool.Version >= r.Version {
				continue
			}

			pool.PardonedAt = r.Now
			if err := pools.Save(ctx, pool, r.Version); err != nil {
				log.WithError(err).Errorln("pools.Save")
				return err
			}
		}

		return nil
	}
}
