package gem

import (
	"context"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(d *db.DB) error {
		tx := d.Update().Model(core.Gem{})

		if err := tx.AutoMigrate(core.Gem{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.GemStore {
	return &gemStore{db: db}
}

type gemStore struct {
	db *db.DB
}

func (s *gemStore) Find(_ context.Context, id string) (*core.Gem, error) {
	gem := core.Gem{ID: id}
	if err := s.db.View().Where("id = ?", id).Take(&gem).Error; err != nil && !db.IsErrorNotFound(err) {
		return nil, err
	}

	return &gem, nil
}

func (s *gemStore) Save(_ context.Context, gem *core.Gem, version int64) error {
	if gem.Version >= version {
		return nil
	}

	if gem.Version == 0 {
		gem.Version = version
		return s.db.Update().Where("id = ?", gem.ID).FirstOrCreate(gem).Error
	}

	updates := map[string]interface{}{
		"amount":    gem.Amount,
		"reward":    gem.Reward,
		"liquidity": gem.Liquidity,
		"profit":    gem.Profit,
		"version":   version,
	}

	tx := s.db.Update().Model(gem).Where("id = ? AND version = ?", gem.ID, gem.Version).Updates(updates)
	if err := tx.Error; err != nil {
		return err
	}

	if tx.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}

	return nil
}

func (s *gemStore) UpdateInfo(_ context.Context, gem *core.Gem) error {
	updates := map[string]interface{}{
		"name":  gem.Name,
		"logo":  gem.Logo,
		"price": gem.Price,
	}

	return s.db.Update().Model(gem).Where("id = ?", gem.ID).Updates(updates).Error
}

func (s *gemStore) List(_ context.Context) ([]*core.Gem, error) {
	var gems []*core.Gem
	if err := s.db.View().Find(&gems).Error; err != nil {
		return nil, err
	}

	return gems, nil
}
