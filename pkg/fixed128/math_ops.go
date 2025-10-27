package fixed128

import (
	"math/big"
)

// Pre-compute 2^64 as a constant for fixed-point operations
// Using a pointer since big.Int is typically passed by reference
// and this avoids copying the 128-bit value
var scale = func() *big.Int {
	scale := big.NewInt(1)
	scale.Lsh(scale, 64) // scale = 2^64
	return scale
}()

// mulFixedPoint performs fixed-point multiplication.
// In fixed-point arithmetic with 64 bits of fractional precision:
// - a represents the value (a_internal / 2^64)
// - b represents the value (b_internal / 2^64)
// - result should be (a_internal * b_internal / 2^64)
func mulFixedPoint(a, b big.Int) big.Int {
	var result big.Int
	result.Mul(&a, &b)
	result.Quo(&result, scale)
	return result
}

// quoFixedPoint performs fixed-point division.
// In fixed-point arithmetic with 64 bits of fractional precision:
// - a represents the value (a_internal / 2^64)
// - b represents the value (b_internal / 2^64)
// - result should be (a_internal * 2^64 / b_internal)
func quoFixedPoint(a, b big.Int) big.Int {
	var result big.Int
	result.Mul(&a, scale)
	result.Quo(&result, &b)
	return result
}
