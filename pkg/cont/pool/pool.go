package pool

import (
	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/holder/pkg/uuid"
	"github.com/fox-one/pkg/logger"
)

func require(condition bool, msg string) error {
	return cont.Require(condition, "Pool/"+msg)
}

func From(r *cont.Request, pools core.PoolStore) (*core.Pool, error) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var id uuid.UUID
	if err := require(r.Scan(&id) == nil, "bad-data"); err != nil {
		return nil, err
	}

	pool, err := pools.Find(ctx, id.String())
	if err != nil {
		log.WithError(err).Errorln("vaults.Find")
		return nil, err
	}

	return pool, nil
}
