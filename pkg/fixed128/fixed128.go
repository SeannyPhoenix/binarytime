// The fixed128 package represents a 128-bit fixed-point fractional number.
// The top 64 bits represent the whole part, and the bottom 64 bits represent
// the fractional part. The underlying data is stored in a big.Int with the divisor

package fixed128

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/bits"
)

type Fixed128 struct {
	value big.Int
}

func NewF128(x int64, y int64) (Fixed128, error) {
	var f128 Fixed128

	if x == math.MinInt64 {
		return f128, fmt.Errorf("value %d is too small to represent in Fixed128", x)
	}

	if y == 0 {
		return f128, fmt.Errorf("division by zero")
	}

	if y == math.MinInt64 {
		return f128, fmt.Errorf("value %d is too small to represent in Fixed128", y)
	}

	negX := x < 0
	var absX uint64
	if negX {
		absX = uint64(-x)
	} else {
		absX = uint64(x)
	}

	negY := y < 0
	var absY uint64
	if negY {
		absY = uint64(-y)
	} else {
		absY = uint64(y)
	}

	hi := getF128Hi(absX, absY)
	lo := getF128Lo(absX, absY)

	var out big.Int
	out.SetUint64(hi)

	var outLow big.Int
	outLow.SetUint64(lo)

	out.Lsh(&out, 64)
	out.Add(&out, &outLow)

	if negX != negY {
		out.Neg(&out)
	}

	f128.value = out
	return f128, nil
}

func (f128 Fixed128) FromF128(y int64) (int64, error) {
	var x int64

	if y == 0 {
		return x, fmt.Errorf("division by zero")
	}

	if y == math.MinInt64 {
		return x, fmt.Errorf("value %d is too small to represent in Fixed128", y)
	}

	negY := y < 0
	var absY uint64
	if negY {
		absY = uint64(-y)
	} else {
		absY = uint64(y)
	}

	hi, lo := hilo(f128)

	x = int64(hi * absY)

	div := absY
	shift := bits.LeadingZeros64(div)
	div <<= shift

	var part uint64
	for i := 0; i < 64 && div > 0; i++ {
		div >>= 1
		bit := lo >> (63 - i) & 1
		part += div * bit
	}
	part >>= shift

	x += int64(part)

	if f128.value.Sign() < 0 != negY {
		x = -x
	}

	return x, nil
}

func (f128 Fixed128) Value() big.Int {
	return f128.value
}

func (f128 Fixed128) String() string {
	hi, lo := hilo(f128)

	var neg rune
	if f128.value.Sign() < 0 {
		neg = '-'
	}

	return fmt.Sprintf("%c%X.%X", neg, hi, lo)
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

func hilo(f128 Fixed128) (uint64, uint64) {
	bytes := f128.value.FillBytes(make([]byte, 16))
	hi := binary.BigEndian.Uint64(bytes[:8])
	lo := binary.BigEndian.Uint64(bytes[8:])
	return hi, lo
}

func (f128 Fixed128) Copy() Fixed128 {
	return Fixed128{
		value: *big.NewInt(0).Set(&f128.value),
	}
}

func (f128 Fixed128) Sign() int {
	return f128.value.Sign()
}

func (f128 Fixed128) Cmp(other *Fixed128) int {
	return f128.value.Cmp(&other.value)
}
