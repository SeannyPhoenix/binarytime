// Package fixed128 provides a 128-bit fixed-point fractional number type.
//
// A Fixed128 represents a signed fixed-point number with 64 bits for the whole
// part and 64 bits for the fractional part, stored internally as a big.Int.
// This allows for precise representation of fractional values without floating-point
// rounding errors, making it suitable for financial calculations, precise arithmetic,
// and other applications requiring exact decimal representation.
//
// # Basic Usage
//
// Create a Fixed128 representing the fraction 22/7 (approximately π):
//
//	f, err := fixed128.New(22, 7)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(f) // Output: 03.249249249249249
//
// Perform arithmetic operations:
//
//	a := fixed128.MustNew(3, 2)  // 1.5
//	b := fixed128.MustNew(5, 4)  // 1.25
//	sum := a.Add(b)              // 2.75
//	product := a.Mul(b)          // 1.875
//
// # Precision and Range
//
// Fixed128 can represent:
//   - Values with whole number part from -2^63 to 2^63-1
//   - Fractional precision of 2^-64 (approximately 5.4×10^-20)
//   - Total numeric range depends on the combination of whole and fractional parts
//
// # Thread Safety
//
// Fixed128 values are immutable after creation. All operations return new
// Fixed128 instances, making them safe for concurrent use without additional
// synchronization.
package fixed128

import (
	"errors"
	"math/big"
)

var (
	// Zero represents the zero value (0.0).
	Zero = Fixed128{}
	// One represents the value 1.0.
	One = func() Fixed128 {
		var one Fixed128
		one.value.SetInt64(1)
		one.value.Lsh(&one.value, 64) // Shift left by 64 bits for fixed-point scaling
		return one
	}()
)

var (
	ErrDivisionByZero         = errors.New("division by zero")
	ErrOverflow               = errors.New("overflow")
	ErrUnderflow              = errors.New("underflow")
	ErrInvalidInput           = errors.New("invalid input")
	ErrMultiplicationOverflow = errors.New("multiplication overflow")
	ErrAdditionOverflow       = errors.New("addition overflow")
)

// Fixed128 represents a 128-bit fixed-point number with 64 bits for the whole part
// and 64 bits for the fractional part. The zero value represents 0.0.
//
// Fixed128 values are immutable; all operations return new instances rather than
// modifying existing ones.
type Fixed128 struct {
	value big.Int
}

// New creates a Fixed128 representing the fraction x/y.
// It returns an error if y is zero (division by zero).
//
// The result represents the exact value x/y with maximum possible precision
// within the Fixed128 format.
//
// Example:
//
//	f, err := fixed128.New(22, 7)  // Creates π approximation (22/7)
//	if err != nil {
//		log.Fatal(err)
//	}
func New(x int64, y int64) (Fixed128, error) {
	return toF128(x, y)
}

// MustNew creates a Fixed128 representing the fraction x/y.
// It panics if y is zero (division by zero).
//
// This is a convenience function for cases where you know y is non-zero
// and want to avoid error handling.
//
// Example:
//
//	half := fixed128.MustNew(1, 2)     // 0.5
//	third := fixed128.MustNew(1, 3)    // 0.333...
func MustNew(x int64, y int64) Fixed128 {
	f128, err := New(x, y)
	if err != nil {
		panic(err)
	}
	return f128
}

// Copy returns a deep copy of the Fixed128.
// Since Fixed128 values are immutable, this is mainly useful for ensuring
// that the internal big.Int is not shared between instances.
func (f128 Fixed128) Copy() Fixed128 {
	return Fixed128{
		value: *big.NewInt(0).Set(&f128.value),
	}
}

// Value returns a copy of the internal big.Int representation.
// This gives access to the raw 128-bit value where the upper 64 bits
// represent the whole part and the lower 64 bits represent the fractional part.
func (f128 Fixed128) Value() big.Int {
	return f128.Copy().value
}

// Sign returns:
//
//	-1 if f128 < 0
//	 0 if f128 == 0
//	+1 if f128 > 0
func (f128 Fixed128) Sign() int {
	return f128.value.Sign()
}

// Cmp compares f128 and other and returns:
//
//	-1 if f128 < other
//	 0 if f128 == other
//	+1 if f128 > other
func (f128 Fixed128) Cmp(other Fixed128) int {
	return f128.value.Cmp(&other.value)
}

// IsNeg reports whether f128 < 0.
func (f128 Fixed128) IsNeg() bool {
	return f128.Sign() < 0
}

// IsZero reports whether f128 == 0.
func (f128 Fixed128) IsZero() bool {
	return f128.value.Sign() == 0
}

// HiLo returns the high and low 64-bit parts of the Fixed128.
// The high part represents the whole number portion,
// and the low part represents the fractional portion.
func (f128 Fixed128) HiLo() (uint64, uint64) {
	_, hi, lo := disassemble(f128)
	return hi, lo
}

// Bytes returns the binary representation of the Fixed128 as a 17-byte slice.
// The first byte indicates the sign (0 for positive, 1 for negative),
// and the remaining 16 bytes contain the absolute value in big-endian format.
func (f128 Fixed128) Bytes() []byte {
	return f128.bytes()
}

// Add returns the sum f128 + b.
func (f128 Fixed128) Add(b Fixed128) Fixed128 {
	var result Fixed128
	result.value.Add(&f128.value, &b.value)
	return result
}

// Sub returns the difference f128 - b.
func (f128 Fixed128) Sub(b Fixed128) Fixed128 {
	var result Fixed128
	result.value.Sub(&f128.value, &b.value)
	return result
}

// Mul returns the product f128 * b.
func (f128 Fixed128) Mul(b Fixed128) Fixed128 {
	return Fixed128{value: mulFixedPoint(f128.value, b.value)}
}

// Quo returns the quotient f128 / b.
// Returns an error if b is zero.
func (f128 Fixed128) Quo(b Fixed128) (Fixed128, error) {
	if b.IsZero() {
		return Fixed128{}, ErrDivisionByZero
	}
	return Fixed128{value: quoFixedPoint(f128.value, b.value)}, nil
}

// MulInt64 multiplies the Fixed128 by an int64 and returns the result as an int64.
// This is useful for converting a fractional Fixed128 back to a whole number
// by multiplying by the original denominator.
//
// Returns an error if y is zero or if the result would overflow int64.
//
// Example:
//
//	f := fixed128.MustNew(22, 7)  // π approximation
//	result, err := f.MulInt64(7)  // Should return 22
func (f128 Fixed128) MulInt64(y int64) (int64, error) {
	return mulInt64(f128, y)
}

// Abs returns the absolute value of the Fixed128.
// For negative values, it returns the value with the sign removed.
// For zero and positive values, it returns the same value.
//
// Example:
//
//	neg := MustNew(-5, 1)
//	pos := neg.Abs()    // 5
//
//	zero := Zero
//	same := zero.Abs()  // 0
func (f128 Fixed128) Abs() Fixed128 {
	var result Fixed128
	result.value.Abs(&f128.value)
	return result
}

// Neg returns the negation (additive inverse) of the Fixed128.
// That is, it returns -f128.
//
// Example:
//
//	pos := MustNew(5, 1)
//	neg := pos.Neg()    // -5
//
//	zero := Zero
//	stillZero := zero.Neg()  // 0 (negation of zero is still zero)
func (f128 Fixed128) Neg() Fixed128 {
	var result Fixed128
	result.value.Neg(&f128.value)
	return result
}
