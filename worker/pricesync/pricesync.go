package pricesync

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/pkg/logger"
)

func New(
	pools core.PoolStore,
	assetz core.AssetService,
) *Syncer {
	return &Syncer{
		pools:  pools,
		assetz: assetz,
	}
}

type Syncer struct {
	pools  core.PoolStore
	assetz core.AssetService
}

func (w *Syncer) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "PriceSync")
	ctx = logger.WithContext(ctx, log)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			_ = w.run(ctx)
		}
	}
}

func (w *Syncer) run(ctx context.Context) error {
	pools, err := w.pools.List(ctx)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("pools.List")
		return err
	}

	for _, pool := range pools {
		asset, err := w.assetz.Find(ctx, pool.ID)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}

			logger.FromContext(ctx).WithError(err).Errorln("assetz.Find")
			continue
		}

		pool.Name = asset.Name
		pool.Logo = asset.Logo
		pool.Price = asset.Price

		if err := w.pools.UpdateInfo(ctx, pool); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("pools.UpdateInfo")
			return err
		}
	}

	return nil
}
