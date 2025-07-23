package binarytime

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/seannyphoenix/binarytime/pkg/byteglyph"
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

func (bt BinaryTime) StringFine() string {
	return fine(bt)
}

func (bt BinaryTime) Glyphs() string {
	bytes := make([]byte, 16)
	bt.value.FillBytes(bytes)

	// return byteglyph.Glyphs(bytes[6:10], 2)
	return byteglyph.Glyphs(bytes, 8)
}

func (bt BinaryTime) TimeGlyphs() string {
	bytes := make([]byte, 16)
	bt.value.FillBytes(bytes)

	return byteglyph.Glyphs(bytes[8:10], 0)
}

func (bt BinaryTime) DateGlyphs() string {
	bytes := make([]byte, 16)
	bt.value.FillBytes(bytes)

	return byteglyph.Glyphs(bytes[6:8], 2)
}

func (bt BinaryTime) DateTimeGlyphs() string {
	bytes := make([]byte, 16)
	bt.value.FillBytes(bytes)

	return byteglyph.Glyphs(bytes[6:10], 2)
}

func (bt BinaryTime) IsZero() bool {
	return bt.value.Sign() == 0
}

func (bt BinaryTime) Value() *big.Int {
	cp := bt.Copy()
	return &cp.value
}

func (bt BinaryTime) Equals(other BinaryTime) bool {
	return bt.value.Cmp(&other.value) == 0
}

const (
	dayNs = uint64(86_400_000_000_000)
	epoch = int64(1 << 32) // 2^32 days before the Unix epoch (1970-01-01T00:00:00Z), or about 11.76 million years
)

func getDays(ns int64) uint64 {
	days := ns / int64(dayNs)
	return uint64(days + epoch)
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

func fine(bt BinaryTime) string {
	bytes := make([]byte, 16)
	bt.value.FillBytes(bytes)

	days := binary.BigEndian.Uint64(bytes[:8])
	subDays := binary.BigEndian.Uint64(bytes[8:])

	return fmt.Sprintf("%016X:%016X", days, subDays)
}
