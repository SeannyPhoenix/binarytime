package fixed128

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/bits"
)

func toF128(x, y int64) (Fixed128, error) {
	if y == 0 {
		return Fixed128{}, fmt.Errorf("division by zero")
	}

	negX, absX := normalize(x)
	negY, absY := normalize(y)

	hi, lo := getComponents(absX, absY)
	f128 := assemble(hi, lo)

	if negX != negY {
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

func getComponents(x, y uint64) (uint64, uint64) {
	if y == 0 {
		panic(fmt.Sprintf("division by zero in getComponents: val %d, div %d", x, y))
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

func assemble(hi, lo uint64) Fixed128 {
	var f128 Fixed128
	f128.value.SetUint64(hi)
	f128.value.Lsh(&f128.value, 64)
	f128.value.Add(&f128.value, big.NewInt(0).SetUint64(lo))
	return f128
}

func fromF128(f128 Fixed128, y int64) (int64, error) {
	var x int64

	if y == 0 {
		return x, fmt.Errorf("division by zero")
	}

	negX := f128.value.Sign() < 0
	negY, absY := normalize(y)

	hi, lo := hilo(f128)

	x = int64(hi * absY)
	part := hydrate(lo, absY)
	x += int64(part)

	if negX != negY {
		x = -x
	}

	return x, nil
}

func hilo(f128 Fixed128) (uint64, uint64) {
	bytes := f128.value.FillBytes(make([]byte, 16))
	hi := binary.BigEndian.Uint64(bytes[:8])
	lo := binary.BigEndian.Uint64(bytes[8:])
	return hi, lo
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

	if shift > 0 {
		part >>= shift - 1
		round := part & 1
		part >>= 1
		part += round
	}

	return part
}
