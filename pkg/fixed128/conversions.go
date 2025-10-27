package fixed128

import (
	"math"
	"math/big"
)

// Error variables for conversions are defined in fixed128.go

// Float64 returns the nearest float64 representation of the Fixed128 value.
//
// Note: This conversion may lose precision due to the limited precision
// of float64 compared to the 128-bit fixed-point representation.
//
// Example:
//
//	f := MustNew(22, 7)  // π approximation
//	val := f.Float64()    // approximately 3.142857142857143
func (f128 Fixed128) Float64() float64 {
	// Convert big.Int to float64
	result := new(big.Float).SetInt(&f128.value)

	// Divide by 2^64 to get the actual fixed-point value
	// We use the scale constant from math_ops.go (which is 2^64)
	scaleFloat := new(big.Float).SetInt(scale)
	result.Quo(result, scaleFloat)

	// Convert to float64
	val, _ := result.Float64()
	return val
}

// FromFloat64 creates a Fixed128 from a float64 value.
// It returns an error if the value is not finite (NaN or Inf).
//
// Note: This conversion has precision limits due to float64's
// 53-bit mantissa. Very large or very small values may lose precision.
//
// Example:
//
//	f, err := FromFloat64(3.14159)
//	if err != nil {
//		log.Fatal(err)
//	}
func FromFloat64(x float64) (Fixed128, error) {
	// Check for NaN or Inf
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return Fixed128{}, ErrInvalidInput
	}

	// Scale the value by 2^64 using big.Float for precision
	val := big.NewFloat(x)
	scaleFloat := new(big.Float).SetInt(scale)
	val.Mul(val, scaleFloat)

	// Convert to big.Int
	bi, _ := val.Int(nil)

	return Fixed128{value: *bi}, nil
}

// Int64 returns the nearest int64 value of the Fixed128.
// The value is rounded to the nearest integer using standard rounding rules.
//
// Returns an error if the value would overflow int64.
//
// Example:
//
//	f := MustNew(7, 2)   // 3.5
//	val, _ := f.Int64()  // returns 4 (rounded)
//
//	g := MustNew(5, 2)   // 2.5
//	val, _ := g.Int64()  // returns 2 (rounded down)
func (f128 Fixed128) Int64() (int64, error) {
	neg, hi, lo := disassemble(f128)

	// Check if we need to round based on the fractional part
	// If lo >= 1<<63, we round up
	shouldRoundUp := lo >= (1 << 63)

	result := hi
	if shouldRoundUp {
		// Check for overflow when rounding up
		if hi == math.MaxInt64 && !neg {
			return 0, ErrOverflow
		}
		if hi == math.MaxUint64 && neg {
			// Would become math.MinInt64
			return math.MinInt64, nil
		}
		result++
	}

	if result > math.MaxInt64 && !neg {
		// Positive overflow
		return 0, ErrOverflow
	}
	if result > uint64(math.MaxInt64)+1 && neg {
		// Negative overflow
		return 0, ErrUnderflow
	}

	intVal := int64(result)
	if neg {
		intVal = -intVal
	}

	return intVal, nil
}

// FromInt64 creates a Fixed128 from an int64 value.
// This is equivalent to creating the fraction x/1.
//
// Example:
//
//	f := FromInt64(42)    // 42.0
//	g := FromInt64(-100)  // -100.0
func FromInt64(x int64) Fixed128 {
	return MustNew(x, 1)
}
