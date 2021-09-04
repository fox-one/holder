package payee

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/cont/gem"
	"github.com/fox-one/holder/pkg/cont/proposal"
	"github.com/fox-one/holder/pkg/cont/sys"
	"github.com/fox-one/holder/pkg/cont/vat"
	"github.com/fox-one/holder/pkg/mtg"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/holder/worker/payee/wallet"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
)

const (
	checkpointKey = "outputs_checkpoint"
)

func New(
	wallets core.WalletStore,
	walletz core.WalletService,
	transactions core.TransactionStore,
	proposals core.ProposalStore,
	properties property.Store,
	gems core.GemStore,
	vaults core.VaultStore,
	parliaments core.Parliament,
	system *core.System,
) *Payee {
	wallets = wallet.BindTransferVersion(wallets)

	actions := map[core.Action]cont.HandlerFunc{
		// sys
		core.ActionSysWithdraw: sys.HandleWithdraw(wallets),
		core.ActionSysProperty: sys.HandleProperty(properties),
		// proposal
		core.ActionProposalMake:  proposal.HandleMake(proposals, walletz, parliaments, system),
		core.ActionProposalShout: proposal.HandleShout(proposals, parliaments, system),
		core.ActionProposalVote:  proposal.HandleVote(proposals, parliaments, walletz, system),
		// gem
		core.ActionGemDonate: gem.HandleDonate(gems, properties),
		core.ActionGemGain:   gem.HandleGain(gems, wallets, system),
		// vat
		core.ActionVaultLock:    vat.HandleLock(gems, vaults),
		core.ActionVaultRelease: vat.HandleRelease(gems, vaults, wallets, properties),
	}

	return &Payee{
		wallets:      wallets,
		properties:   properties,
		transactions: transactions,
		system:       system,
		actions:      actions,
	}
}

type Payee struct {
	wallets      core.WalletStore
	properties   property.Store
	transactions core.TransactionStore
	gems         core.GemStore
	vaults       core.VaultStore
	system       *core.System
	actions      map[core.Action]cont.HandlerFunc
}

func (w *Payee) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "payee")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = 500 * time.Millisecond
			}
		}
	}
}

func (w *Payee) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	v, err := w.properties.Get(ctx, checkpointKey)
	if err != nil {
		log.WithError(err).Errorln("properties.Get", err)
		return err
	}

	const Limit = 500
	outputs, err := w.wallets.List(ctx, v.Int64(), Limit)
	if err != nil {
		log.WithError(err).Errorln("wallets.List")
		return err
	}

	if len(outputs) == 0 {
		return errors.New("EOF")
	}

	for _, u := range outputs {
		if err := w.handleOutput(ctx, u); err != nil {
			return err
		}

		if err := w.properties.Save(ctx, checkpointKey, u.ID); err != nil {
			log.WithError(err).Errorln("properties.Save", checkpointKey)
			return err
		}
	}

	return nil
}

func (w *Payee) handleOutput(ctx context.Context, output *core.Output) error {
	log := logger.FromContext(ctx).WithField("output", output.TraceID)
	ctx = logger.WithContext(ctx, log)

	message := decodeMemo(output.Memo)
	req := requestFromOutput(output)

	// bind system version
	sysVersion, err := w.properties.Get(ctx, sys.SystemVersionKey)
	if err != nil {
		log.WithError(err).Errorln("properties.Get", sys.SystemVersionKey)
		return err
	}

	if v := sysVersion.Int(); v > req.SysVersion {
		req.SysVersion = v
	}

	if payload, err := core.DecodeTransactionAction(message); err == nil {
		if req.Body, err = mtg.Scan(payload.Body, &req.Action); err == nil {
			if follow, _ := uuid.FromBytes(payload.FollowID); follow != uuid.Zero {
				req.FollowID = follow.String()
			}

			return w.handleRequest(req.WithContext(ctx))
		}
	}

	return nil
}

func (w *Payee) handleRequest(r *cont.Request) error {
	ctx := r.Context()
	log := logger.FromContext(ctx).WithField("action", r.Action.String())

	h, ok := w.actions[r.Action]
	if !ok {
		log.Debugf("handler not found")
		return nil
	}

	tx := r.Tx()

	if err := h(r); err != nil {
		var e cont.Error
		if !errors.As(err, &e) {
			return err
		}

		if r.Sender != "" && cont.ShouldRefund(e.Flag) {
			memo := core.TransferAction{
				ID:     r.FollowID,
				Source: e.Error(),
			}.Encode()

			transfer := &core.Transfer{
				TraceID:   uuid.Modify(r.TraceID, memo),
				AssetID:   r.AssetID,
				Amount:    r.Amount,
				Memo:      memo,
				Threshold: 1,
				Opponents: []string{r.Sender},
			}

			if err := w.wallets.CreateTransfers(ctx, []*core.Transfer{transfer}); err != nil {
				log.WithError(err).Errorln("wallets.CreateTransfers")
				return err
			}
		}

		tx.Status = core.TransactionStatusAbort
		tx.Message = e.Msg
	} else {
		tx.Status = core.TransactionStatusOk
	}

	tx.Parameters, _ = json.Marshal(r.Values())
	if err := w.transactions.Create(ctx, tx); err != nil {
		log.WithError(err).Errorln("transactions.Create")
		return err
	}

	if r.Next != nil {
		return w.handleRequest(r.Next)
	}

	return nil
}

func decodeMemo(memo string) []byte {
	if b, err := base64.StdEncoding.DecodeString(memo); err == nil {
		return b
	}

	if b, err := base64.URLEncoding.DecodeString(memo); err == nil {
		return b
	}

	return []byte(memo)
}

func requestFromOutput(output *core.Output) *cont.Request {
	return &cont.Request{
		Now:        output.CreatedAt,
		Version:    output.ID,
		SysVersion: 1,
		TraceID:    output.TraceID,
		Sender:     output.Sender,
		FollowID:   output.TraceID,
		AssetID:    output.AssetID,
		Amount:     output.Amount,
	}
}
