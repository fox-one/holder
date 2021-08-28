package gem

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/mtg/types"
	"github.com/fox-one/pkg/logger"
)

func HandleDonate(gems core.GemStore) cont.HandlerFunc {
	return func(r *cont.Request) error {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		gem, err := From(r.WithBody(types.UUID(r.AssetID)), gems)
		if err != nil {
			return err
		}

		if gem.Version < r.Version {
			gem.Amount = gem.Amount.Add(r.Amount)
			if err := gems.Save(ctx, gem, r.Version); err != nil {
				log.WithError(err).Errorln("gems.Save")
				return err
			}
		}

		return nil
	}
}
