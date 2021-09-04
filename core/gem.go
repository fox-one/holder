package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Gem struct {
		ID        string          `sql:"size:36;primary_key" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		UpdatedAt time.Time       `json:"updated_at,omitempty"`
		Version   int64           `json:"version,omitempty"`
		Name      string          `sql:"size:64" json:"name,omitempty"`
		Logo      string          `sql:"size:256" json:"logo,omitempty"`
		Amount    decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Reward    decimal.Decimal `sql:"type:decimal(64,8)" json:"reward,omitempty"`
		Liquidity decimal.Decimal `sql:"type:decimal(64,12)" json:"liquidity,omitempty"`
		Profit    decimal.Decimal `sql:"type:decimal(64,8)" json:"profit,omitempty"`
		Price     decimal.Decimal `sql:"type:decimal(24,8)" json:"price,omitempty"`
	}

	GemStore interface {
		Find(ctx context.Context, id string) (*Gem, error)
		Save(ctx context.Context, gem *Gem, version int64) error
		UpdateInfo(ctx context.Context, gem *Gem) error
		List(ctx context.Context) ([]*Gem, error)
	}
)
