package fixed128

import (
	"errors"
	"math/bits"
)

var (
	ErrorDivisionByZero       = errors.New("division by zero")
	ErrorAdditionOverflow     = errors.New("addition overflow")
	ErrorSubtractionUnderflow = errors.New("subtraction underflow")
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
