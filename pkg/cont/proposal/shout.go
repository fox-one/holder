package proposal

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/pkg/logger"
)

func HandleShout(
	proposals core.ProposalStore,
	parliaments core.Parliament,
	system *core.System,
) cont.HandlerFunc {
	return func(r *cont.Request) error {
		ctx := r.Context()

		if err := require(system.IsMember(r.Sender), "not-member"); err != nil {
			return err
		}

		p, err := From(r, proposals)
		if err != nil {
			return err
		}

		if err := parliaments.ProposalCreated(ctx, p); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("parliaments.ProposalCreated")
			return err
		}

		return nil
	}
}
