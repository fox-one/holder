package proposal

import (
	"encoding/base64"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/mtg"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string) error {
	return cont.Require(condition, "proposal/"+msg)
}

func From(r *cont.Request, proposals core.ProposalStore) (*core.Proposal, error) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	p, err := proposals.Find(ctx, id.String())
	if err != nil {
		log.WithError(err).Errorln("proposals.Find")
		return nil, err
	}

	if err := require(p.ID > 0, "not init"); err != nil {
		return nil, err
	}

	return p, nil
}

func handleProposal(r *cont.Request, walletz core.WalletService, system *core.System, action core.Action, p *core.Proposal) error {
	pid, _ := uuid.FromString(p.TraceID)
	data, _ := mtg.Encode(action, pid)
	data, _ = core.TransactionAction{
		Body: data,
	}.Encode()
	memo := base64.StdEncoding.EncodeToString(data)

	ctx := r.Context()
	if err := walletz.HandleTransfer(ctx, &core.Transfer{
		TraceID:   uuid.Modify(r.TraceID, p.TraceID+system.ClientID),
		AssetID:   system.GasAssetID,
		Amount:    system.GasAmount,
		Threshold: system.Threshold,
		Opponents: system.Members,
		Memo:      memo,
	}); err != nil {
		logger.FromContext(ctx).WithError(err).Errorf("handle proposal %s failed", action.String())
		return err
	}

	return nil
}
