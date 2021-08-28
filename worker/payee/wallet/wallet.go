package wallet

import (
	"context"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/cont"
)

func BindTransferVersion(wallets core.WalletStore) core.WalletStore {
	return &withVersion{wallets}
}

type withVersion struct {
	core.WalletStore
}

func (s *withVersion) CreateTransfers(ctx context.Context, transfers []*core.Transfer) error {
	version := cont.VersionFrom(ctx)

	for _, transfer := range transfers {
		transfer.Version = version
	}

	return s.WalletStore.CreateTransfers(ctx, transfers)
}
