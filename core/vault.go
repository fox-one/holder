package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type VaultStatus int

const (
	VaultStatusLocking VaultStatus = iota
	VaultStatusReleased
)

//go:generate stringer -type VaultStatus -trimprefix VaultStatus

type (
	Vault struct {
		ID          int64           `sql:"primary_key" json:"id,omitempty"`
		CreatedAt   time.Time       `json:"created_at,omitempty"`
		UpdatedAt   time.Time       `json:"updated_at,omitempty"`
		ReleasedAt  time.Time       `json:"released_at,omitempty"`
		Version     int64           `json:"version,omitempty"`
		TraceID     string          `sql:"size:36" json:"trace_id,omitempty"`
		UserID      string          `sql:"size:36" json:"user_id,omitempty"`
		Status      VaultStatus     `json:"status,omitempty"`
		AssetID     string          `sql:"size:36" json:"asset_id,omitempty"`
		Duration    int64           `json:"duration,omitempty"`
		MinDuration int64           `json:"min_duration,omitempty"`
		Amount      decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Liquidity   decimal.Decimal `sql:"type:decimal(64,8)" json:"liquidity,omitempty"`
		Reward      decimal.Decimal `sql:"type:decimal(64,8)" json:"reward,omitempty"`
		Penalty     decimal.Decimal `sql:"type:decimal(64,8)" json:"penalty,omitempty"`
	}

	VaultStore interface {
		Create(ctx context.Context, vault *Vault) error
		Update(ctx context.Context, vault *Vault, version int64) error
		Find(ctx context.Context, traceID string) (*Vault, error)
		List(ctx context.Context, fromID int64, limit int) ([]*Vault, error)
		ListUser(ctx context.Context, userID string) ([]*Vault, error)
	}
)
