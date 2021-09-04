package pricesync

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/pkg/logger"
)

func New(
	gems core.GemStore,
	assetz core.AssetService,
) *Syncer {
	return &Syncer{
		gems:   gems,
		assetz: assetz,
	}
}

type Syncer struct {
	gems   core.GemStore
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
	gems, err := w.gems.List(ctx)
	if err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("gems.List")
		return err
	}

	for _, gem := range gems {
		asset, err := w.assetz.Find(ctx, gem.ID)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}

			logger.FromContext(ctx).WithError(err).Errorln("assetz.Find")
			continue
		}

		gem.Name = asset.Name
		gem.Logo = asset.Logo
		gem.Price = asset.Price

		if err := w.gems.UpdateInfo(ctx, gem); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("gems.UpdateInfo")
			return err
		}
	}

	return nil
}
