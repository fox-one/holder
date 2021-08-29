package main

import (
	"github.com/fox-one/holder/parliament"
	"github.com/google/wire"
)

var parliamentSet = wire.NewSet(
	parliament.New,
)
