package core

import (
	"github.com/asaskevich/govalidator"
	"github.com/shopspring/decimal"
)

// System stores system information.
type System struct {
	Admins       []string
	ClientID     string
	ClientSecret string
	Members      []string
	Threshold    uint8
	GasAssetID   string
	GasAmount    decimal.Decimal
	Version      string
}

func (s *System) IsMember(id string) bool {
	return govalidator.IsIn(id, s.Members...)
}

func (s *System) IsStaff(id string) bool {
	return govalidator.IsIn(id, s.Admins...)
}
