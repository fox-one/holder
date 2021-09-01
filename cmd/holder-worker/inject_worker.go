package main

import (
	"github.com/fox-one/holder/cmd/holder-worker/config"
	"github.com/fox-one/holder/worker"
	"github.com/fox-one/holder/worker/assigner"
	"github.com/fox-one/holder/worker/cashier"
	"github.com/fox-one/holder/worker/datadog"
	"github.com/fox-one/holder/worker/events"
	"github.com/fox-one/holder/worker/messenger"
	"github.com/fox-one/holder/worker/payee"
	"github.com/fox-one/holder/worker/spentsync"
	"github.com/fox-one/holder/worker/syncer"
	"github.com/fox-one/holder/worker/txsender"
	"github.com/google/wire"
)

var workerSet = wire.NewSet(
	provideCashierConfig,
	cashier.New,
	messenger.New,
	payee.New,
	spentsync.New,
	syncer.New,
	txsender.New,
	assigner.New,
	provideDataDogConfig,
	datadog.New,
	events.New,
	provideWorkers,
)

func provideCashierConfig() cashier.Config {
	return cashier.Config{
		Batch:    _flag.cashier.batch,
		Capacity: _flag.cashier.capacity,
	}
}

func provideDataDogConfig(cfg *config.Config) datadog.Config {
	return datadog.Config{
		ConversationID: cfg.DataDog.ConversationID,
		Interval:       _flag.datadog.interval,
		Version:        version,
	}
}

func provideWorkers(
	a *cashier.Cashier,
	b *messenger.Messenger,
	c *payee.Payee,
	d *events.Events,
	e *spentsync.SpentSync,
	f *txsender.Sender,
	g *syncer.Syncer,
	h *assigner.Assigner,
	i *datadog.Datadog,
) []worker.Worker {
	workers := []worker.Worker{a, b, c, d, e, f, g, h, i}
	return workers
}
