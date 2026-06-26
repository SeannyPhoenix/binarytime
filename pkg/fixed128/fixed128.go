package fixed128

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
	return add(f128, other)
}

// Sub subtracts another Fixed128 number from the current one and returns the result.
// It returns an error if the subtraction results in an underflow.
func (f128 Fixed128) Sub(other Fixed128) (Fixed128, error) {
	return sub(f128, other)
}

func (f128 Fixed128) Cmp(other Fixed128) int {
	if f128.neg && !other.neg {
		return -1
	}
	if !f128.neg && other.neg {
		return 1
	}

	if f128.hi < other.hi {
		if f128.neg {
			return 1
		}
		return -1
	}
	if f128.hi > other.hi {
		if f128.neg {
			return -1
		}
		return 1
	}

	if f128.lo < other.lo {
		if f128.neg {
			return 1
		}
		return -1
	}
	if f128.lo > other.lo {
		if f128.neg {
			return -1
		}
		return 1
	}

	return 0
}
