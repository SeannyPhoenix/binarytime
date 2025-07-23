package binarytime

import (
	"encoding/binary"
	"fmt"
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
	nanos := t.UnixNano()

	// handle negative times later
	if nanos < 0 {
		return BinaryTime{}
	}

	return FromNanos(nanos)
}

func FromNanos(nanos int64) BinaryTime {
	days := getDays(nanos)
	subDays := getSubDay(uint64(nanos))

	upper := toBytes(days)
	lower := toBytes(subDays)
	bytes := append(upper[:], lower[:]...)

	bt := new(big.Int).SetBytes(bytes)

	return BinaryTime{
		value: *bt,
	}
}

func (bt BinaryTime) Copy() BinaryTime {
	return BinaryTime{
		value: *big.NewInt(0).Set(&bt.value),
	}
}

func (bt BinaryTime) String() string {
	return coarse(bt)
}

func (bt BinaryTime) IsZero() bool {
	return bt.value.Sign() == 0
}

func (bt BinaryTime) Value() *big.Int {
	cp := bt.Copy()
	return &cp.value
}

const dayNs = uint64(86_400_000_000_000)

func getDays(ns int64) uint64 {
	if ns < 0 {
		ns = -ns
	}
	return uint64(ns) / dayNs
}

func getSubDay(ns uint64) uint64 {
	ns %= dayNs
	d := dayNs

	var sd uint64
	var i int
	for ; i < 64 && d > 1; i++ {
		d >>= 1

		v := ns / d
		if v == 1 {
			ns -= d
		}

		sd <<= 1
		sd |= v
	}

	sd <<= (64 - i)
	return sd
}

func toBytes(v uint64) [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], v)
	return b
}

func coarse(bt BinaryTime) string {
	bytes := make([]byte, 16)
	bt.value.FillBytes(bytes)

	days := binary.BigEndian.Uint16(bytes[6:8])
	subDays := binary.BigEndian.Uint16(bytes[8:10])

	return fmt.Sprintf("%04X:%04X", days, subDays)
}
