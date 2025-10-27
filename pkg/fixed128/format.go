package fixed128

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// String implements fmt.Stringer.
// It returns a hex representation of the Fixed128 value.
func (f128 Fixed128) String() string {
	if f128.IsZero() {
		return "00.00"
	}

	b := f128.bytes()

	high := 1
	for high < 8 && b[high] == 0 {
		high++
	}

	low := 17
	for low > 10 && b[low-1] == 0 {
		low--
	}

	out, err := stringWithPrecision(b, high, low)
	if err != nil {
		panic(err)
	}

	return out
}

// StringWithPrecision returns a string representation with custom precision control.
// The high parameter controls how many leading bytes to show for the whole part,
// and low controls how many trailing bytes to show for the fractional part.
//
// Returns an error if the precision parameters are invalid.
func (f128 Fixed128) StringWithPrecision(high int, low int) (string, error) {
	return stringWithPrecision(f128.bytes(), high, low)
}

func stringWithPrecision(b []byte, high int, low int) (string, error) {
	if len(b) != 17 {
		return "", fmt.Errorf("invalid length: %d", len(b))
	}
	if high >= 9 || low <= 9 {
		return "", fmt.Errorf("invalid precision: %d, %d", high, low)
	}

	var s strings.Builder
	if b[0] == 1 {
		s.WriteRune('-')
	}
	s.WriteString(hex.EncodeToString(b[high:9]))
	s.WriteRune('.')
	s.WriteString(hex.EncodeToString(b[9:low]))

	return s.String(), nil
}

// Parse parses a string representation of a Fixed128.
// The string should be in the format produced by String(), e.g., "03.14159" or "-01.50".
//
// Returns an error if the string is not in the correct format.
func Parse(s string) (Fixed128, error) {
	var f128 Fixed128
	if err := f128.UnmarshalText([]byte(s)); err != nil {
		return Fixed128{}, err
	}
	return f128, nil
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
