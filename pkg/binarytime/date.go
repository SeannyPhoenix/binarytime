package binarytime

import (
	"time"

	"github.com/seannyphoenix/binarytime/pkg/fixed128"
)

var (
	BinaryTimeOffset = fixed128.FromParts(1<<42, 0, false)
)

type Date struct {
	value fixed128.Fixed128
}

func Now() Date {
	return DateFromTime(time.Now())
}

func DateFromTime(t time.Time) Date {
	return DateFromUnixNanos(t.UnixNano())
}

// DateFromUnixNanos creates a BinaryTime from a Unix timestamp in nanoseconds.
func DateFromUnixNanos(nanos int64) Date {
	value, err := fixed128.ByDivision(nanos, dayNs)
	if err != nil {
		return Date{}
	}

	value, err = value.Add(BinaryTimeOffset)
	if err != nil {
		return Date{}
	}

	return Date{value: value}
}

func (d Date) Time() time.Time {
	return time.Unix(0, d.UnixNano())
}

func (d Date) UnixNano() int64 {
	v, err := d.value.Sub(BinaryTimeOffset)
	if err != nil {
		return 0
	}
	ns, _ := v.MulInt64(dayNs)
	return ns
}

// IsZero checks if the BinaryTime is zero.
func (d Date) IsZero() bool {
	return d.value.IsZero()
}

// Equals checks if two BinaryTime instances are equal.
func (d Date) Equals(other Date) bool {
	return d.value.Cmp(other.value) == 0
}

// Fixed128 returns the underlying Fixed128 value of the Date.
// This is a copy of the value, not a reference.
func (d Date) Fixed128() fixed128.Fixed128 {
	return d.value
}

func (d Date) Bytes() []byte {
	return d.value.Bytes()
}

func DateFromBytes(b []byte) (Date, error) {
	value, err := fixed128.FromBytes(b)
	d := Date{value: value}
	if err != nil {
		return d, err
	}
	return d, nil
}
