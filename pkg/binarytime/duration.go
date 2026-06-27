package binarytime

import (
	"time"

	"github.com/seannyphoenix/binarytime/pkg/fixed128"
)

type Duration struct {
	value fixed128.Fixed128
}

func FromDuration(d time.Duration) Duration {
	return FromNanos(d.Nanoseconds())
}

func FromNanos(nanos int64) Duration {
	v, err := fixed128.ByDivision(nanos, dayNs)
	if err != nil {
		return Duration{}
	}
	return Duration{value: v}
}
