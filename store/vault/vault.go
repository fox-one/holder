package vault

import (
	"context"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(d *db.DB) error {
		tx := d.Update().Model(core.Vault{})

		if err := tx.AutoMigrate(core.Vault{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_vaults_trace", "trace_id").Error; err != nil {
			return err
		}

		if err := tx.AddIndex("idx_vaults_user", "user_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.VaultStore {
	return &vaultStore{db: db}
}

type vaultStore struct {
	db *db.DB
}

func (s *vaultStore) Create(_ context.Context, vault *core.Vault) error {
	return s.db.Update().Where("trace_id = ?", vault.TraceID).FirstOrCreate(vault).Error
}

func (s *vaultStore) Update(_ context.Context, vault *core.Vault, version int64) error {
	if vault.Version >= version {
		return nil
	}

	updates := map[string]interface{}{
		"released_at": vault.ReleasedAt,
		"status":      vault.Status,
		"reward":      vault.Reward,
		"penalty":     vault.Penalty,
		"version":     version,
	}

	tx := s.db.Update().Model(vault).Where("version = ?", vault.Version).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}

	return nil
}

func (s *vaultStore) Find(_ context.Context, traceID string) (*core.Vault, error) {
	vault := core.Vault{TraceID: traceID}
	if err := s.db.View().Where("trace_id = ?", traceID).Take(&vault).Error; err != nil && !db.IsErrorNotFound(err) {
		return nil, err
	}

	return &vault, nil
}

func (s *vaultStore) List(ctx context.Context, fromID int64, limit int) ([]*core.Vault, error) {
	var vaults []*core.Vault
	if err := s.db.View().Where("id > ?", fromID).Order("id").Limit(limit).Find(&vaults).Error; err != nil {
		return nil, err
	}

	return vaults, nil
}

func (s *vaultStore) ListUser(_ context.Context, userID string) ([]*core.Vault, error) {
	var vaults []*core.Vault
	if err := s.db.View().Where("user_id = ?", userID).Order("id DESC").Find(&vaults).Error; err != nil {
		return nil, err
	}

	return vaults, nil
}
