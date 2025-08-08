// The fixed128 package represents a 128-bit fixed-point fractional number.
// The top 64 bits represent the whole part, and the bottom 64 bits represent
// the fractional part. The underlying data is stored in a big.Int with the divisor

package fixed128

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/bits"
)

type Fixed128 struct {
	value   big.Int
	divisor uint64
}

func NewF128(val int64, div uint64) (Fixed128, error) {
	return toF128(val, div)
}

func toF128(val int64, div uint64) (Fixed128, error) {
	f128 := Fixed128{divisor: div}

	if div == 0 {
		return f128, fmt.Errorf("division by zero")
	}

	neg := val < 0
	if neg {
		val = -val
	}

	abs := uint64(val)
	hi := getF128Hi(abs, div)
	lo := getF128Lo(abs, div)

	var out big.Int
	out.SetUint64(hi)
	out.Lsh(&out, 64)

	var outLow big.Int
	outLow.SetUint64(lo)

	out.Add(&out, &outLow)

	if neg {
		out.Neg(&out)
	}

	f128.value = out
	return f128, nil
}

func getF128Hi(val, div uint64) uint64 {
	return val / div
}

func getF128Lo(val, div uint64) uint64 {
	if div == 0 {
		return 0
	}

	part := val % div

	shift := bits.LeadingZeros64(div)
	div <<= shift
	part <<= shift

	var out uint64
	var i int
	for ; i < 64 && div > 1 && part > 0; i++ {
		div >>= 1
		bit := part / div
		part -= bit * div
		out <<= 1
		out |= bit
	}

	out <<= (64 - i)
	return out
}

func FromF128(f128 Fixed128) (int64, error) {
	val, _, err := fromF128(f128)
	return val, err
}

func fromF128(f128 Fixed128) (int64, uint64, error) {
	bytes := f128.value.FillBytes(make([]byte, 16))
	hi := binary.BigEndian.Uint64(bytes[:8])
	lo := binary.BigEndian.Uint64(bytes[8:])

	d := f128.divisor
	full := int64(hi * d)

	for lo > 0 {
		d >>= 1
		lo = bits.RotateLeft64(lo, 1)
		bit := lo & 1
		full += int64(d * bit)
		lo >>= 1
		lo <<= 1
	}

	if f128.value.Sign() < 0 {
		full = -full
	}

	return full, f128.divisor, nil
}
