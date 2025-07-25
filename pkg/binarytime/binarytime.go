package binarytime

import (
	"encoding/binary"
	"math/big"
	"time"
)

type BinaryTime struct {
	value big.Int
}

func Now() BinaryTime {
	return FromTime(time.Now())
}

func FromTime(t time.Time) BinaryTime {
	return FromUnixNanos(t.UnixNano())
}

// FromUnixNanos creates a BinaryTime from a Unix timestamp in nanoseconds.
func FromUnixNanos(nanos int64) BinaryTime {
	days := getDays(nanos)
	subDays := getSubDay(nanos)

	var upper, lower [8]byte
	binary.BigEndian.PutUint64(upper[:], days)
	binary.BigEndian.PutUint64(lower[:], subDays)

	bytes := append(upper[:], lower[:]...)
	bt := new(big.Int).SetBytes(bytes)

	return BinaryTime{
		value: *bt,
	}
}

// Copy creates a copy of the BinaryTime.
// This is useful for ensuring that the original BinaryTime is not modified.
func (bt BinaryTime) Copy() BinaryTime {
	return BinaryTime{
		value: *big.NewInt(0).Set(&bt.value),
	}
}

// IsZero checks if the BinaryTime is zero.
func (bt BinaryTime) IsZero() bool {
	return bt.value.Sign() == 0
}

// Equals checks if two BinaryTime instances are equal.
func (bt BinaryTime) Equals(other BinaryTime) bool {
	return bt.value.Cmp(&other.value) == 0
}

// Value returns the underlying big.Int value of the BinaryTime.
// This is acopy of the value, not a reference.
func (bt BinaryTime) Value() big.Int {
	return bt.value
}

const (
	// There are 86.4 trillion nanoseconds in a day.
	dayNs = int64(86_400_000_000_000)

	// 2^32 plus 2^48 days before the Unix epoch (1970-01-01T00:00:00Z),
	// or ~770.64 billion years + ~11.76 million years
	epoch = int64(1<<32 + 1<<48)
)

// Because we're adding the epoch to the days, we know that the days will always be positive.
// This means that we can safely cast the days to uint64 without worrying about negative values.
func getDays(ns int64) uint64 {
	days := ns / dayNs
	days += epoch
	if ns < 0 {
		days -= 1
	}

	return uint64(days)
}

func getSubDay(ns int64) uint64 {
	ns %= dayNs
	if ns < 0 {
		ns = (ns + dayNs) % dayNs
	}

	// Now that we've taken care of the negative case,
	// we can safely cast ns to uint64.
	sub := uint64(ns)
	d := uint64(dayNs)
	var sd uint64
	var i int
	for ; i < 64 && d > 1; i++ {
		d >>= 1

		v := sub / d
		if v == 1 {
			sub -= d
		}

		sd <<= 1
		sd |= v
	}

	sd <<= (64 - i)
	return sd
}
