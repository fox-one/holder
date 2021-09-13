package pool

import (
	"context"

	"github.com/fox-one/holder/core"
)

func Cache(pools core.PoolStore) core.PoolStore {
	return &cachePools{
		PoolStore: pools,
		c:         make(map[string]*core.Pool, 12),
	}
}

type cachePools struct {
	core.PoolStore
	c map[string]*core.Pool
}

func (s *cachePools) Find(ctx context.Context, id string) (*core.Pool, error) {
	if p, ok := s.c[id]; ok {
		return p, nil
	}

	pool, err := s.PoolStore.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	s.c[id] = pool
	return pool, nil
}
