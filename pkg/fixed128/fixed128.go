// The fixed128 package represents a 128-bit fixed-point fractional number.
// The top 64 bits represent the whole part, and the bottom 64 bits represent
// the fractional part. The underlying data is stored in a big.Int with the divisor

package fixed128

import (
	"fmt"
	"math/big"
)

type Fixed128 struct {
	value big.Int
}

func NewF128(x int64, y int64) (Fixed128, error) {
	return toF128(x, y)
}

func (f128 Fixed128) FromF128(y int64) (int64, error) {
	return fromF128(f128, y)
}

func (f128 Fixed128) String() string {
	hi, lo := hilo(f128)

	var neg rune
	if f128.Sign() < 0 {
		neg = '-'
	}

	return fmt.Sprintf("%c%X.%X", neg, hi, lo)
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
