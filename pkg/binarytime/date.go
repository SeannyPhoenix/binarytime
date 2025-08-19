package binarytime

import (
	"math/big"
	"time"

	"github.com/seannyphoenix/binarytime/pkg/fixed128"
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
	value, err := fixed128.NewF128(nanos, dayNs)
	if err != nil {
		return Date{}
	}

	return Date{value: value}
}

func (bt Date) ToTime() time.Time {
	ns, err := bt.value.FromF128(dayNs)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(0, ns)
}

// Copy creates a copy of the BinaryTime.
// This is useful for ensuring that the original BinaryTime is not modified.
func (bt Date) Copy() Date {
	return Date{
		value: bt.value.Copy(),
	}
}

// IsZero checks if the BinaryTime is zero.
func (bt Date) IsZero() bool {
	return bt.value.Sign() == 0
}

// Equals checks if two BinaryTime instances are equal.
func (bt Date) Equals(other Date) bool {
	return bt.value.Cmp(other.value) == 0
}

// Value returns the underlying Fixed128 value of the BinaryTime.
// This is a copy of the value, not a reference.
func (bt Date) Fixed128() fixed128.Fixed128 {
	f128 := bt.value
	return f128
}

func (bt Date) BigInt() big.Int {
	return bt.value.Value()
}

func (bt Date) Bytes() []byte {
	bytes := make([]byte, 16)
	bi := bt.BigInt()
	bi.FillBytes(bytes[:])
	return bytes
}
