package sys

import (
	"strings"

	"github.com/fox-one/holder/pkg/cont"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
)

const (
	SystemVersionKey    = "system_version"
	SystemProfitRateKey = "system_profit_rate"
	SystemDonateFeeRate = "system_donate_fee_rate"
)

func HandleProperty(
	properties property.Store,
) cont.HandlerFunc {
	return func(r *cont.Request) error {
		ctx := r.Context()
		if err := require(r.Gov(), "not-authorized"); err != nil {
			return err
		}

		var key, value string
		if err := require(r.Scan(&key, &value) == nil, "bad-data"); err != nil {
			return err
		}

		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if err := require(key != "" && value != "", "empty"); err != nil {
			return err
		}

		if err := properties.Save(ctx, key, value); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("properties.Save", key, value)
			return err
		}

		return nil
	}
}
