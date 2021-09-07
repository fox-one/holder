package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Pool struct {
		ID        string          `sql:"size:36;primary_key" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		UpdatedAt time.Time       `json:"updated_at,omitempty"`
		Version   int64           `json:"version,omitempty"`
		Amount    decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Share     decimal.Decimal `sql:"type:decimal(64,12)" json:"share,omitempty"`
		Reward    decimal.Decimal `sql:"type:decimal(64,8)" json:"reward,omitempty"`
		Liquidity decimal.Decimal `sql:"type:decimal(64,12)" json:"liquidity,omitempty"`
		Profit    decimal.Decimal `sql:"type:decimal(64,8)" json:"profit,omitempty"`
		// Pool info
		Name  string          `sql:"size:64" json:"name,omitempty"`
		Logo  string          `sql:"size:256" json:"logo,omitempty"`
		Price decimal.Decimal `sql:"type:decimal(24,8)" json:"price,omitempty"`
	}

	PoolStore interface {
		Find(ctx context.Context, id string) (*Pool, error)
		Save(ctx context.Context, pool *Pool, version int64) error
		UpdateInfo(ctx context.Context, pool *Pool) error
		List(ctx context.Context) ([]*Pool, error)
	}
)
