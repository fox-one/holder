package pool

import (
	"context"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(d *db.DB) error {
		tx := d.Update().Model(core.Pool{})

		if err := tx.AutoMigrate(core.Pool{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.PoolStore {
	return &poolStore{db: db}
}

type poolStore struct {
	db *db.DB
}

func (s *poolStore) Find(_ context.Context, id string) (*core.Pool, error) {
	pool := core.Pool{ID: id}
	if err := s.db.View().Where("id = ?", id).Take(&pool).Error; err != nil && !db.IsErrorNotFound(err) {
		return nil, err
	}

	return &pool, nil
}

func (s *poolStore) Save(_ context.Context, pool *core.Pool, version int64) error {
	if pool.Version >= version {
		return nil
	}

	if pool.Version == 0 {
		pool.Version = version
		return s.db.Update().Where("id = ?", pool.ID).FirstOrCreate(pool).Error
	}

	updates := map[string]interface{}{
		"amount":    pool.Amount,
		"share":     pool.Share,
		"reward":    pool.Reward,
		"liquidity": pool.Liquidity,
		"profit":    pool.Profit,
		"version":   version,
	}

	tx := s.db.Update().Model(pool).Where("id = ? AND version = ?", pool.ID, pool.Version).Updates(updates)
	if err := tx.Error; err != nil {
		return err
	}

	if tx.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}

	return nil
}

func (s *poolStore) UpdateInfo(_ context.Context, pool *core.Pool) error {
	updates := map[string]interface{}{
		"name":  pool.Name,
		"logo":  pool.Logo,
		"price": pool.Price,
	}

	return s.db.Update().Model(pool).Where("id = ?", pool.ID).Updates(updates).Error
}

func (s *poolStore) List(_ context.Context) ([]*core.Pool, error) {
	var pools []*core.Pool
	if err := s.db.View().Find(&pools).Error; err != nil {
		return nil, err
	}

	return pools, nil
}
