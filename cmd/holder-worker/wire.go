//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pandodao/bholdings/cmd/bholdings-worker/config"
)

func buildApp(cfg *config.Config) (app, error) {
	wire.Build(
		storeSet,
		serviceSet,
		notifierSet,
		parliamentSet,
		workerSet,
		serverSet,
		wire.Struct(new(app), "*"),
	)

	return app{}, nil
}
