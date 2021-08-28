package parliament

import (
	"context"

	"github.com/fox-one/holder/pkg/uuid"
)

func (s *parliament) fetchAssetSymbol(ctx context.Context, assetID string) string {
	if uuid.IsNil(assetID) {
		return "ALL"
	}

	coin, err := s.assetz.Find(ctx, assetID)
	if err != nil {
		return "NULL"
	}

	return coin.Symbol
}

func (s *parliament) fetchUserName(ctx context.Context, userID string) string {
	user, err := s.userz.Find(ctx, userID)
	if err != nil {
		return "NULL"
	}

	return user.Name
}
