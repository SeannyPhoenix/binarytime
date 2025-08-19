package fixed128

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
)

func toF128(x, y int64) (Fixed128, error) {
	if y == 0 {
		return Fixed128{}, fmt.Errorf("division by zero")
	}

	negX, absX := normalize(x)
	negY, absY := normalize(y)
	neg := negX != negY

	hi, lo := getComponents(absX, absY)

	f128 := assemble(neg, hi, lo)
	return f128, nil
}

func normalize(v int64) (bool, uint64) {
	mask := uint64(v >> 63)
	neg := mask != 0
	abs := (uint64(v) ^ mask) - mask
	return neg, abs
}

func getComponents(x, y uint64) (uint64, uint64) {
	if y == 0 {
		panic(fmt.Sprintf("division by zero in getComponents: x %d, y %d", x, y))
	}

	var hi, lo uint64
	hi = x / y
	part := x % y

	shift := bits.LeadingZeros64(y)
	y <<= shift
	part <<= shift

	var i int
	for ; i < 64 && y > 1 && part > 0; i++ {
		y >>= 1
		bit := part / y
		part -= bit * y
		lo <<= 1
		lo |= bit
	}

	lo <<= (64 - i)

	return hi, lo
}

func assemble(neg bool, hi, lo uint64) Fixed128 {
	var f128 Fixed128

	var buf [16]byte
	binary.BigEndian.PutUint64(buf[:8], hi)
	binary.BigEndian.PutUint64(buf[8:], lo)
	f128.value.SetBytes(buf[:])

	if neg {
		f128.value.Neg(&f128.value)
	}

	return f128
}

func disassemble(f128 Fixed128) (bool, uint64, uint64) {
	var buf [16]byte
	f128.value.FillBytes(buf[:])
	hi := binary.BigEndian.Uint64(buf[:8])
	lo := binary.BigEndian.Uint64(buf[8:])

	neg := f128.IsNeg()

	return neg, hi, lo
}

func fromF128(f128 Fixed128, y int64) (int64, error) {
	if y == 0 {
		return 0, fmt.Errorf("division by zero")
	}

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
