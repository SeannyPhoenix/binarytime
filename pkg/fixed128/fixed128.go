package fixed128

import "encoding/binary"

var (
	Zero        = Fixed128{}
	One         = Fixed128{hi: 1}
	NegativeOne = Fixed128{hi: 1, neg: true}
)

// Fixed128 represents a 128-bit fixed-point fractional number.
// The top 64 bits represent the whole part, and the bottom 64 bits represent
// the fractional part. The underlying data is stored in two uint64 values,
// with a boolean flag to indicate if the number is negative.
type Fixed128 struct {
	hi  uint64
	lo  uint64
	neg bool
}

// FromParts creates a Fixed128 from its constituent parts:
// hi: the high 64 bits (whole part),
// lo: the low 64 bits (fractional part),
// neg: a boolean indicating if the number is negative.
func FromParts(hi, lo uint64, neg bool) Fixed128 {
	return Fixed128{hi: hi, lo: lo, neg: neg}
}

// ByDivision divides two int64 numbers and returns the result as a Fixed128.
// It returns an error if the divisor is zero.
func ByDivision(x, y int64) (Fixed128, error) {
	return divide(x, y)
}

func MustByDivision(x, y int64) Fixed128 {
	f128, err := ByDivision(x, y)
	if err != nil {
		panic(err)
	}
	return f128
}

func (f128 Fixed128) IsZero() bool {
	return f128 == Zero
}

// Parts returns the constituent parts of the Fixed128 number:
// hi: the high 64 bits (whole part),
// lo: the low 64 bits (fractional part),
// neg: a boolean indicating if the number is negative.
func (f128 Fixed128) Parts() (uint64, uint64, bool) {
	return f128.hi, f128.lo, f128.neg
}

// Add adds two Fixed128 numbers and returns the result.
// It returns an error if the addition results in an overflow.
func (f128 Fixed128) Add(other Fixed128) (Fixed128, error) {
	// Can safely add same-signed values
	if f128.neg == other.neg {
		return add(f128, other)
	}

	// If the signs are different, we need to perform subtraction
	// Determine which value is larger in absolute terms
	cmp := absCmp(f128, other)
	switch cmp {
	case 0:
		// If they are equal, the result is zero
		return Zero, nil
	case 1:
		// f128 is larger in absolute terms
		result, err := sub(f128, other)
		if err != nil {
			return Zero, err
		}
		result.neg = f128.neg // Result takes the sign of the larger value
		return result, nil
	case -1:
		// other is larger in absolute terms
		result, err := sub(other, f128)
		if err != nil {
			return Zero, err
		}
		result.neg = other.neg // Result takes the sign of the larger value
		return result, nil
	default:
		return Zero, ErrorAdditionOverflow // This should never happen
	}
}

// Sub subtracts another Fixed128 number from the current one and returns the result.
// It returns an error if the subtraction results in an underflow.
func (f128 Fixed128) Sub(other Fixed128) (Fixed128, error) {
	return f128.Add(other.Negate())
}

func (f128 Fixed128) Negate() Fixed128 {
	return Fixed128{
		hi:  f128.hi,
		lo:  f128.lo,
		neg: !f128.neg,
	}
}

func (f128 Fixed128) Cmp(other Fixed128) int {
	if f128.neg && !other.neg {
		return -1
	}
	if !f128.neg && other.neg {
		return 1
	}

	cmpResult := absCmp(f128, other)
	if f128.neg {
		cmpResult = -cmpResult
	}
	return cmpResult
}

func (f128 Fixed128) Bytes() []byte {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[:8], f128.hi)
	binary.BigEndian.PutUint64(b[8:], f128.lo)
	return b
}

func FromBytes(b []byte) (Fixed128, error) {
	f128 := Fixed128{}
	if len(b) != 16 {
		return f128, ErrorBadByteLength
	}
	f128.hi = binary.BigEndian.Uint64(b[:8])
	f128.lo = binary.BigEndian.Uint64(b[8:])
	return f128, nil
}

func (f128 Fixed128) Sign() bool {
	return f128.neg
}

// MulInt64 multiplies the Fixed128 by an int64 and returns the result as int64.
// It returns an error if the multiplication would overflow.
func (f128 Fixed128) MulInt64(multiplier int64) (int64, error) {
	return mulInt64(f128, multiplier)
}
