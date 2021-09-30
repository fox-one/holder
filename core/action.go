package core

import (
	"encoding/base64"
	"encoding/json"

	"github.com/fox-one/holder/pkg/mtg/types"
	"github.com/fox-one/msgpack"
)

//go:generate stringer -type Action -trimprefix Action

type Action types.BitInt

const (
	_ Action = iota
	ActionSysWithdraw
	ActionSysProperty
	ActionProposalMake
	ActionProposalShout
	ActionProposalVote
	ActionPoolDonate
	ActionPoolGain
	ActionVaultLock
	ActionVaultRelease
	ActionPoolPardon
	ActionPoolPardonAll
)

func (i Action) MarshalBinary() (data []byte, err error) {
	return types.BitInt(i).MarshalBinary()
}

func (i *Action) UnmarshalBinary(data []byte) error {
	var b types.BitInt
	if err := b.UnmarshalBinary(data); err != nil {
		return err
	}

	*i = Action(b)
	return nil
}

type TransactionAction struct {
	FollowID []byte `msgpack:"f,omitempty"`
	Body     []byte `msgpack:"b,omitempty"`
}

func (action TransactionAction) Encode() ([]byte, error) {
	return msgpack.Marshal(action)
}

func DecodeTransactionAction(b []byte) (*TransactionAction, error) {
	var action TransactionAction
	if err := msgpack.Unmarshal(b, &action); err != nil {
		return nil, err
	}

	return &action, nil
}

const (
	TransferSourceRefund = "Refund"
)

type TransferAction struct {
	ID     string `json:"id,omitempty"`
	Source string `json:"s,omitempty"`
}

func (action TransferAction) Encode() string {
	b, _ := json.Marshal(action)
	return base64.StdEncoding.EncodeToString(b)
}
