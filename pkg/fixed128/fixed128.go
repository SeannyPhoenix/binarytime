// The fixed128 package represents a 128-bit fixed-point fractional number.
// The top 64 bits represent the whole part, and the bottom 64 bits represent
// the fractional part. The underlying data is stored in a big.Int.

package fixed128

import (
	"math/big"
)

type Fixed128 struct {
	value big.Int
}

func NewF128(x int64, y int64) (Fixed128, error) {
	return toF128(x, y)
}

func MustNewF128(x int64, y int64) Fixed128 {
	f128, err := NewF128(x, y)
	if err != nil {
		panic(err)
	}
	return f128
}

func (f128 Fixed128) FromF128(y int64) (int64, error) {
	return fromF128(f128, y)
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

func (f128 Fixed128) HiLo() (uint64, uint64) {
	_, hi, lo := disassemble(f128)
	return hi, lo
}

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

func (f128 Fixed128) IsNeg() bool {
	return f128.Sign() < 0
}

func (f128 Fixed128) IsZero() bool {
	return f128.value.Sign() == 0
}

func (f128 Fixed128) Bytes() []byte {
	return f128.bytes()
}
