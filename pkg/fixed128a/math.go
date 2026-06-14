package fixed128a

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/bits"
)

func toF128(x, y int64) (Fixed128, error) {
	if y == 0 {
		return Fixed128{}, fmt.Errorf("division by zero")
	}

	negX, absX := normalize(x)
	negY, absY := normalize(y)
	neg := negX != negY

	xBig := big.NewInt(0).SetUint64(absX)
	yBig := big.NewInt(0).SetUint64(absY)

	f128 := Fixed128{value: *xBig}.Lsh(64)
	f128.value.Div(&f128.value, yBig)

	if neg {
		f128.value.Neg(&f128.value)
	}

	return f128, nil
}

func normalize(v int64) (bool, uint64) {
	mask := uint64(v >> 63)
	neg := mask != 0
	abs := (uint64(v) ^ mask) - mask
	return neg, abs
}

func disassemble(f128 Fixed128) (bool, uint64, uint64) {
	var buf [16]byte
	f128.value.FillBytes(buf[:])
	hi := binary.BigEndian.Uint64(buf[:8])
	lo := binary.BigEndian.Uint64(buf[8:])

	neg := f128.IsNeg()

	return neg, hi, lo
}

func mulInt64(f128 Fixed128, y int64) (int64, error) {
	negX, hi, lo := disassemble(f128)
	negY, absY := normalize(y)

	whole, err := multiply64(hi, absY)
	if err != nil {
		return 0, err
	}

	part := hydrate(lo, absY)
	x, err := add64(whole, part)
	if err != nil {
		return 0, fmt.Errorf("addition overflow: %d + %d", whole, part)
	}

	if negX != negY {
		x = -x
	}

	return x, nil
}

func multiply64(a, b uint64) (uint64, error) {
	hi, lo := bits.Mul64(a, b)
	if hi > 0 {
		return 0, fmt.Errorf("multiplication overflow: %d * %d", a, b)
	}

	return lo, nil
}

func add64(a, b uint64) (int64, error) {
	sum, carry := bits.Add64(a, b, 0)
	if sum > math.MaxInt64 || carry > 0 {
		return 0, fmt.Errorf("addition overflow: %d + %d", a, b)
	}
	return int64(sum), nil
}

func hydrate(lo, div uint64) uint64 {
	shift := bits.LeadingZeros64(div)
	div <<= shift

	var part uint64
	for i := 0; i < 64 && div > 0; i++ {
		div >>= 1
		bit := lo >> (63 - i) & 1
		part += div * bit
	}

	part = round(shift, part)
	return part
}

func round(shift int, part uint64) uint64 {
	if shift == 0 {
		return part
	}

	part >>= shift - 1
	bit := part & 1
	part >>= 1
	part += bit
	return part
}
