// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/fox-one/holder/cmd/holder-worker/config"
	"github.com/fox-one/holder/handler/node"
	"github.com/fox-one/holder/parliament"
	"github.com/fox-one/holder/service/asset"
	message2 "github.com/fox-one/holder/service/message"
	"github.com/fox-one/holder/service/user"
	wallet2 "github.com/fox-one/holder/service/wallet"
	"github.com/fox-one/holder/store/message"
	"github.com/fox-one/holder/store/pool"
	"github.com/fox-one/holder/store/proposal"
	"github.com/fox-one/holder/store/transaction"
	user2 "github.com/fox-one/holder/store/user"
	"github.com/fox-one/holder/store/vault"
	"github.com/fox-one/holder/store/wallet"
	"github.com/fox-one/holder/worker/assigner"
	"github.com/fox-one/holder/worker/cashier"
	"github.com/fox-one/holder/worker/datadog"
	"github.com/fox-one/holder/worker/events"
	"github.com/fox-one/holder/worker/keeper"
	"github.com/fox-one/holder/worker/messenger"
	"github.com/fox-one/holder/worker/payee"
	"github.com/fox-one/holder/worker/pricesync"
	"github.com/fox-one/holder/worker/spentsync"
	"github.com/fox-one/holder/worker/syncer"
	"github.com/fox-one/holder/worker/txsender"
	"github.com/fox-one/pkg/store/property"
)

// Injectors from wire.go:

func buildApp(cfg *config.Config) (app, error) {
	db, err := provideDatabase(cfg)
	if err != nil {
		return app{}, err
	}
	walletStore := wallet.New(db)
	client, err := provideMixinClient(cfg)
	if err != nil {
		return app{}, err
	}
	walletConfig := provideWalletServiceConfig(cfg)
	walletService := wallet2.New(client, walletConfig)
	system := provideSystem(cfg)
	cashierConfig := provideCashierConfig()
	cashierCashier := cashier.New(walletStore, walletService, system, cashierConfig)
	messageStore := message.New(db)
	messageService := message2.New(client)
	messengerMessenger := messenger.New(messageStore, messageService)
	transactionStore := transaction.New(db)
	proposalStore := proposal.New(db)
	store := propertystore.New(db)
	poolStore := pool.New(db)
	vaultStore := vault.New(db)
	userConfig := _wireConfigValue
	userService := user.New(client, userConfig)
	assetService := asset.New(client)
	coreParliament := parliament.New(messageStore, userService, assetService, walletService, system)
	payeePayee := payee.New(walletStore, walletService, transactionStore, proposalStore, store, poolStore, vaultStore, coreParliament, assetService, system)
	userStore := user2.New(db)
	localizer, err := provideLocalizer(cfg)
	if err != nil {
		return app{}, err
	}
	notifierConfig := provideNotifierConfig(cfg)
	notifier := provideNotifier(system, assetService, messageStore, poolStore, vaultStore, userStore, walletService, localizer, notifierConfig)
	eventsEvents := events.New(transactionStore, notifier, store)
	spentSync := spentsync.New(walletStore, notifier)
	sender := txsender.New(walletStore)
	syncerSyncer := syncer.New(walletStore, walletService, store)
	assignerAssigner := assigner.New(walletStore, system)
	datadogConfig := provideDataDogConfig(cfg)
	datadogDatadog := datadog.New(walletStore, store, messageService, datadogConfig)
	keeperKeeper := keeper.New(poolStore, vaultStore, notifier, walletService, system)
	pricesyncSyncer := pricesync.New(poolStore, assetService)
	v := provideWorkers(cashierCashier, messengerMessenger, payeePayee, eventsEvents, spentSync, sender, syncerSyncer, assignerAssigner, datadogDatadog, keeperKeeper, pricesyncSyncer)
	server := node.New(system, store)
	mux := provideRoute(server)
	serverServer := provideServer(mux)
	mainApp := app{
		workers: v,
		server:  serverServer,
	}
	return mainApp, nil
}

var (
	_wireConfigValue = user.Config{}
)
