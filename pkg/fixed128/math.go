package fixed128

import (
	"errors"
	"math/bits"
)

var (
	ErrorDivisionByZero       = errors.New("division by zero")
	ErrorAdditionOverflow     = errors.New("addition overflow")
	ErrorSubtractionUnderflow = errors.New("subtraction underflow")
	ErrorBadByteLength        = errors.New("bad byte length")
)

func divide(x, y int64) (Fixed128, error) {
	if y == 0 {
		return Zero, ErrorDivisionByZero
	}

	negX, absX := normalize(x)
	negY, absY := normalize(y)
	neg := negX != negY

	hi, lo := getParts(absX, absY)

	return Fixed128{hi: hi, lo: lo, neg: neg}, nil
}

func normalize(v int64) (bool, uint64) {
	mask := uint64(v >> 63)
	neg := mask != 0
	abs := (uint64(v) ^ mask) - mask
	return neg, abs
}

func getParts(x, y uint64) (uint64, uint64) {
	var hi, lo uint64

	hi = x / y
	part := x % y

	shift := bits.LeadingZeros64(y)
	y <<= shift
	part <<= shift

	var i int
	for ; i < 64 && y > 0 && part > 0; i++ {
		y >>= 1
		bit := 1 & ^(uint64(int64(part-y) >> 63))
		part -= bit * y
		lo <<= 1
		lo |= bit
	}

	lo <<= (64 - i)
	return hi, lo
}

func add(a, b Fixed128) (Fixed128, error) {
	lo, carry := bits.Add64(a.lo, b.lo, 0)
	hi, carry := bits.Add64(a.hi, b.hi, carry)
	if carry != 0 {
		return Zero, ErrorAdditionOverflow
	}

	return Fixed128{
		hi:  hi,
		lo:  lo,
		neg: a.neg,
	}, nil
}

func sub(a, b Fixed128) (Fixed128, error) {
	lo, borrow := bits.Sub64(a.lo, b.lo, 0)
	hi, borrow := bits.Sub64(a.hi, b.hi, borrow)
	if borrow != 0 {
		return Zero, ErrorSubtractionUnderflow
	}

	return Fixed128{
		hi:  hi,
		lo:  lo,
		neg: a.neg,
	}, nil
}

func absCmp(a, b Fixed128) int {
	hiDiff := int64(a.hi - b.hi)
	loDiff := int64(a.lo - b.lo)

	hiSign := hiDiff>>63 - ((-hiDiff) >> 63)
	loSign := loDiff>>63 - ((-loDiff) >> 63)

	mask := (hiDiff | -hiDiff) >> 63

	return int((hiSign & mask) | (loSign & ^mask))
}

func mulInt64(f128 Fixed128, multiplier int64) (int64, error) {
	// Handle sign
	negMul := multiplier < 0
	absMul := uint64(multiplier)
	if negMul {
		absMul = uint64(-multiplier)
	}

	// Multiply hi and lo parts using 64-bit multiplication
	hiHigh, hiLow := bits.Mul64(f128.hi, absMul)
	loHigh, _ := bits.Mul64(f128.lo, absMul)

	// Add fractional overflow from lo to hi result
	result, carry := bits.Add64(hiLow, loHigh, 0)
	if hiHigh != 0 || carry != 0 {
		return 0, ErrorAdditionOverflow
	}

	// Apply sign: flip if f128 and multiplier have different signs
	signResult := int64(result)
	if f128.neg != negMul {
		signResult = -signResult
	}

	return signResult, nil
}
