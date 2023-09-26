package keeper

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/mtg"
	"github.com/fox-one/holder/pkg/mtg/types"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/holder/worker/keeper/pool"
	"github.com/fox-one/pkg/logger"
)

func New(
	pools core.PoolStore,
	vaults core.VaultStore,
	notifier core.Notifier,
	walletz core.WalletService,
	system *core.System,
) *Keeper {
	return &Keeper{
		pools:    pools,
		vaults:   vaults,
		notifier: notifier,
		walletz:  walletz,
		system:   system,
		filter:   make(map[int64]struct{}),
	}
}

type Keeper struct {
	pools    core.PoolStore
	vaults   core.VaultStore
	notifier core.Notifier
	system   *core.System
	walletz  core.WalletService
	filter   map[int64]struct{}
}

func (w *Keeper) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "keeper")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	interval := 10 * time.Minute

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-time.After(dur):
			if err := w.run(ctx, t); err != nil {
				dur = time.Second
			} else {
				dur = t.Truncate(interval).Add(interval).Sub(t)
			}
		}
	}
}

func (w *Keeper) run(ctx context.Context, t time.Time) error {
	var (
		log   = logger.FromContext(ctx)
		pools = pool.Cache(w.pools)
	)

	var fromID int64 = 0
	const limit = 100

	for {
		vaults, err := w.vaults.List(ctx, fromID, limit)
		if err != nil {
			log.WithError(err).Errorln("vaults.List")
			return err
		}

		for _, vault := range vaults {
			fromID = vault.ID

			if _, ok := w.filter[vault.ID]; ok {
				continue
			}

			if vault.Status == core.VaultStatusReleased {
				continue
			}

			pool, err := pools.Find(ctx, vault.AssetID)
			if err != nil {
				log.WithError(err).Errorln("pools.Find")
				return err
			}

			pool.Reform(vault)
			if vault.EndAt().After(t) {
				continue
			}

			if err := w.notifier.LockDone(ctx, pool, vault); err != nil {
				log.WithError(err).Errorln("notifier.LockDone")
				return err
			}

			if err := w.releaseVault(ctx, vault); err != nil {
				return err
			}

			w.filter[vault.ID] = struct{}{}
		}

		if len(vaults) < limit {
			break
		}
	}

	return nil
}

func (w *Keeper) releaseVault(ctx context.Context, vault *core.Vault) error {
	body, err := mtg.Encode(core.ActionVaultRelease, types.UUID(vault.TraceID))
	if err != nil {
		return err
	}

	action := core.TransactionAction{
		Body: body,
	}

	b, err := action.Encode()
	if err != nil {
		return err
	}

	t := &core.Transfer{
		TraceID:   uuid.Modify(vault.TraceID, "keep:release"),
		AssetID:   w.system.GasAssetID,
		Amount:    w.system.GasAmount,
		Memo:      base64.StdEncoding.EncodeToString(b),
		Threshold: w.system.Threshold,
		Opponents: w.system.Members,
	}

	if err := w.walletz.HandleTransfer(ctx, t); err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("walletz.HandleTransfer")
		return err
	}

	return nil
}
