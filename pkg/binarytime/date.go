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
	value, err := fixed128.New(nanos, dayNs)
	if err != nil {
		return Date{}
	}

	return Date{value: value}
}

func (d Date) ToTime() time.Time {
	ns, err := d.value.MulInt64(dayNs)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(0, ns)
}

func (d Date) Base64() string {
	return d.value.Base64()
}

// Copy creates a copy of the BinaryTime.
// This is useful for ensuring that the original BinaryTime is not modified.
func (d Date) Copy() Date {
	return Date{
		value: d.value.Copy(),
	}
}

// IsZero checks if the BinaryTime is zero.
func (d Date) IsZero() bool {
	return d.value.Sign() == 0
}

// Equals checks if two BinaryTime instances are equal.
func (d Date) Equals(other Date) bool {
	return d.value.Cmp(other.value) == 0
}

// Value returns the underlying Fixed128 value of the BinaryTime.
// This is a copy of the value, not a reference.
func (d Date) Fixed128() fixed128.Fixed128 {
	f128 := d.value
	return f128
}

func (d Date) BigInt() big.Int {
	return d.value.Value()
}

func (d Date) Bytes() []byte {
	bytes := make([]byte, 16)
	bi := d.BigInt()
	bi.FillBytes(bytes[:])
	return bytes
}
