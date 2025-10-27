package fixed128

import (
	"math"
	"math/big"
	"strings"
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

// DecimalString returns a decimal string representation of the Fixed128 value.
// The result is in human-readable decimal format, e.g. "3.14159" instead of
// the hex format returned by String().
//
// Example:
//
//	f := MustNew(22, 7)  // π approximation
//	dec := f.DecimalString()  // "3.142857142857143"
func (f128 Fixed128) DecimalString() string {
	return f128.DecimalStringWithPrecision(15)
}

// DecimalStringWithPrecision returns a decimal string representation with
// the specified number of decimal places.
//
// The precision parameter determines how many decimal places to show.
// If a value requires more precision than requested, it will be rounded.
//
// This implementation uses big.Float to maintain full 128-bit precision,
// avoiding the limitations of float64.
//
// Example:
//
//	f := MustNew(1, 3)  // 0.333...
//	dec := f.DecimalStringWithPrecision(5)  // "0.33333"
func (f128 Fixed128) DecimalStringWithPrecision(precision int) string {
	if f128.IsZero() {
		if precision == 0 {
			return "0"
		}
		return "0." + strings.Repeat("0", precision)
	}

	// Use big.Float for maximum precision (maintains 128-bit precision)
	// Work with absolute value to avoid sign duplication
	var absVal big.Int
	absVal.Abs(&f128.value)

	val := new(big.Float).SetInt(&absVal)
	scaleFloat := new(big.Float).SetInt(scale)
	val.Quo(val, scaleFloat)

	// Set high precision to avoid rounding errors
	val.SetPrec(256)

	// Handle sign and build result
	var result strings.Builder
	neg := f128.IsNeg()

	// For precision 0, use banker's rounding (round half to even)
	if precision == 0 {
		intPart, _ := val.Int(nil)

		// Check if we need to round by getting the fractional part
		intFloat := new(big.Float).SetInt(intPart)
		frac := new(big.Float).Sub(val, intFloat)

		half := big.NewFloat(0.5)
		half.SetPrec(256)

		cmp := frac.Cmp(half)
		if cmp > 0 {
			// > 0.5, round up
			intPart.Add(intPart, big.NewInt(1))
		} else if cmp == 0 {
			// Exactly 0.5, round to nearest even (banker's rounding)
			if intPart.Bit(0) == 1 {
				// Odd number, round away from zero
				intPart.Add(intPart, big.NewInt(1))
			}
			// Even number, already correct (truncate)
		}
		// < 0.5, already correct (truncate)

		if neg {
			result.WriteByte('-')
		}
		result.WriteString(intPart.String())
		return result.String()
	}

	// For decimal places, format and add sign if needed
	str := val.Text('f', precision)
	if neg {
		result.WriteByte('-')
	}
	result.WriteString(str)

	return result.String()
}
