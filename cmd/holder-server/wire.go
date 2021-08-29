//+build wireinject

package main

import (
	"github.com/fox-one/holder/cmd/holder-server/config"
	"github.com/fox-one/holder/server"
	"github.com/google/wire"
)

func buildServer(cfg *config.Config) (*server.Server, error) {
	wire.Build(
		storeSet,
		serviceSet,
		notifierSet,
		sessionSet,
		serverSet,
	)

	return &server.Server{}, nil
}
