package fixed128legacy

import (
	"fmt"
	"math/big"
)

var (
	ErrorDivisionByZero = fmt.Errorf("division by zero")
)

type Fixed128 struct {
	value big.Int
}

func New(x, y int64) (Fixed128, error) {
	var f128 Fixed128

	if y == 0 {
		return f128, ErrorDivisionByZero
	}

	absX, negX := normalize(x)
	absY, negY := normalize(y)
	neg := negX != negY

	xBig := big.NewInt(0).SetUint64(absX)
	yBig := big.NewInt(0).SetUint64(absY)

	f128.value.Lsh(xBig, 64)
	f128.value.Div(&f128.value, yBig)

	if neg {
		f128.value.Neg(&f128.value)
	}

	return f128, nil
}

func MustNew(x, y int64) Fixed128 {
	f128, err := New(x, y)
	if err != nil {
		panic(err)
	}
	return f128
}

func normalize(v int64) (uint64, bool) {
	mask := uint64(v >> 63)
	neg := mask != 0
	abs := (uint64(v) ^ mask) - mask
	return abs, neg
}

func (f128 Fixed128) Copy() Fixed128 {
	return Fixed128{
		value: *big.NewInt(0).Set(&f128.value),
	}
}

func (f128 Fixed128) Value() big.Int {
	return f128.Copy().value
}

func (f128 Fixed128) Sign() int {
	return f128.value.Sign()
}

func (f128 Fixed128) Cmp(other Fixed128) int {
	return f128.value.Cmp(&other.value)
}

func (f128 Fixed128) IsNeg() bool {
	return f128.Sign() < 0
}

func (f128 Fixed128) IsZero() bool {
	return f128.value.Sign() == 0
}

// func (f128 Fixed128) HiLo() (uint64, uint64) {
// 	_, hi, lo := disassemble(f128)
// 	return hi, lo
// }

// func (f128 Fixed128) Bytes() []byte {
// 	return f128.bytes()
// }

func (f128 Fixed128) Add(b Fixed128) Fixed128 {
	var result Fixed128
	result.value.Add(&f128.value, &b.value)
	return result
}

func (f128 Fixed128) Sub(b Fixed128) Fixed128 {
	var result Fixed128
	result.value.Sub(&f128.value, &b.value)
	return result
}

func (f128 Fixed128) Mul(b Fixed128) Fixed128 {
	var result Fixed128
	result.value.Mul(&f128.value, &b.value)
	return result
}

func (f128 Fixed128) Quo(b Fixed128) (Fixed128, error) {
	var result Fixed128
	result.value.Quo(&f128.value, &b.value)
	return result, nil
}

func (f128 Fixed128) Lsh(bits uint) Fixed128 {
	var result Fixed128
	result.value.Lsh(&f128.value, bits)
	return result
}

// func (f128 Fixed128) MulInt64(y int64) (int64, error) {
// 	return mulInt64(f128, y)
// }
