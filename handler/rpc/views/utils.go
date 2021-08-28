package views

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Milliseconds(t *time.Time) float64 {
	if t == nil {
		return 0
	}

	return float64(t.UnixNano() / int64(time.Millisecond))
}

func Unix(t *time.Time) float64 {
	if t == nil {
		return 0
	}

	return float64(t.Unix())
}

func Time(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}
