package sys

import (
	"github.com/fox-one/holder/pkg/cont"
)

func require(condition bool, msg string) error {
	return cont.Require(condition, "Sys/"+msg)
}
