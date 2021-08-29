package main

import (
	"github.com/fox-one/holder/cmd/holder-server/config"
	"github.com/fox-one/holder/store/gem"
	"github.com/fox-one/holder/store/message"
	"github.com/fox-one/holder/store/proposal"
	"github.com/fox-one/holder/store/transaction"
	"github.com/fox-one/holder/store/user"
	"github.com/fox-one/holder/store/vault"
	"github.com/fox-one/holder/store/wallet"
	"github.com/fox-one/pkg/store/db"
	propertystore "github.com/fox-one/pkg/store/property"
	"github.com/google/wire"
)

var storeSet = wire.NewSet(
	provideDatabase,
	proposal.New,
	transaction.New,
	user.New,
	vault.New,
	wallet.New,
	message.New,
	propertystore.New,
	gem.New,
)

func provideDatabase(cfg *config.Config) (*db.DB, error) {
	database, err := db.Open(cfg.DB)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(database); err != nil {
		return nil, err
	}

	return database, nil
}
