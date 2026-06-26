package binarytime

import (
	"time"

	"github.com/seannyphoenix/binarytime/pkg/uint128"
)

type Duration struct {
	value uint128.uint128
}

func FromDuration(d time.Duration) Duration {
	return FromNanoseconds(d.Nanoseconds())
}

func FromNanoseconds(nanos int64) Duration {
	v, err := uint128.New(nanos, dayNs)
	if err != nil {
		return Duration{}
	}
	return Duration{value: v}
}
