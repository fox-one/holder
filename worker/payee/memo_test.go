package payee

import (
	"testing"

	"github.com/fox-one/holder/core"
	"github.com/fox-one/holder/pkg/mtg"
	"github.com/fox-one/holder/pkg/mtg/types"
)

func TestDecodeMemo(t *testing.T) {
	const memo = "gqFmxBBorDHcALNLYqqBS/EA9sHDoWLELQEDAQEQM21dlzKcMw2OYit8m6QOoAYAAAAAAgEQgBfSAHhwS4K1P3S64dLa1w=="

	b := decodeMemo(memo)
	payload, err := core.DecodeTransactionAction(b)
	if err != nil {
		t.Error(err)
		return
	}

	var action types.BitInt
	mtg.Scan(payload.Body, &action)
	t.Log(action)
}
