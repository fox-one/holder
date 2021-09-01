package keeper

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/mtg"
	"github.com/fox-one/holder/pkg/mtg/types"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func New(
	vaults core.VaultStore,
	walletz core.WalletService,
	system *core.System,
) *Keeper {
	return &Keeper{
		vaults:  vaults,
		walletz: walletz,
		system:  system,
	}
}

type Keeper struct {
	vaults  core.VaultStore
	walletz core.WalletService
	system  *core.System
}

func (w *Keeper) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "keeper")
	ctx = logger.WithContext(ctx, log)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-time.After(5 * time.Minute):
			_ = w.run(ctx, t)
		}
	}
}

func (w *Keeper) run(ctx context.Context, t time.Time) error {
	traceID := uuid.MD5(t.Format(time.RFC3339))

	if err := w.ping(ctx, traceID); err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("ping")
		return err
	}

	if err := w.handleExpiredVaults(ctx, traceID); err != nil {
		logger.FromContext(ctx).WithError(err).Errorln("handleExpiredVaults")
		return err
	}

	return nil
}

// ping 会发送一条转账到多签组，刷新 utxo 同步的 checkpoint
func (w *Keeper) ping(ctx context.Context, traceID string) error {
	return w.walletz.HandleTransfer(ctx, &core.Transfer{
		TraceID:   uuid.Modify(traceID, "ping"),
		AssetID:   w.system.GasAssetID,
		Amount:    w.system.GasAmount,
		Memo:      "ping",
		Threshold: w.system.Threshold,
		Opponents: w.system.Members,
	})
}

func (w *Keeper) handleExpiredVaults(ctx context.Context, traceID string) error {
	log := logger.FromContext(ctx)

	var fromID int64 = 0
	const limit = 100

	for {
		vaults, err := w.vaults.List(ctx, fromID, limit)
		if err != nil {
			log.WithError(err).Errorln("vaults.List")
			return err
		}

		for _, v := range vaults {
			fromID = v.ID

			if v.Status != core.VaultStatusLocking {
				continue
			}

			if time.Since(v.CreatedAt).Milliseconds()/1000 < v.Duration {
				continue
			}

			memo := buildMemo(core.ActionVaultRelease, types.UUID(v.TraceID))
			if err := w.walletz.HandleTransfer(ctx, &core.Transfer{
				TraceID:   uuid.Modify(traceID, v.TraceID),
				AssetID:   w.system.GasAssetID,
				Amount:    w.system.GasAmount,
				Threshold: w.system.Threshold,
				Opponents: w.system.Members,
				Memo:      memo,
			}); err != nil {
				log.WithError(err).Errorln("walletz.HandleTransfer")
				return err
			}
		}

		if len(vaults) < limit {
			break
		}
	}

	return nil
}

func buildMemo(values ...interface{}) string {
	body, err := mtg.Encode(values...)
	if err != nil {
		panic(err)
	}

	action := core.TransactionAction{
		Body: body,
	}

	data, err := action.Encode()
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(data)
}
