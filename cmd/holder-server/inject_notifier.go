package main

import (
	"github.com/fox-one/holder/notifier"
	"github.com/google/wire"
)

var notifierSet = wire.NewSet(
	wire.Value(notifier.Config{}),
	notifier.New,
)
