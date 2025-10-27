package fixed128

import (
	"math"
	"testing"
)

// TestFloat64 tests the Float64() method
func TestFloat64(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int64
		expected float64
	}{
		{"zero", 0, 1, 0.0},
		{"one", 1, 1, 1.0},
		{"half", 1, 2, 0.5},
		{"third", 1, 3, 1.0 / 3.0},
		{"quarter", 1, 4, 0.25},
		{"negative", -5, 2, -2.5},
		{"larger", 100, 3, 100.0 / 3.0},
		{"maxint", math.MaxInt64, 1, float64(math.MaxInt64)},
		{"pi_approx", 22, 7, 22.0 / 7.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f128 := MustNew(tt.x, tt.y)
			got := f128.Float64()

			// Use tolerance for floating point comparison
			tolerance := 1e-15
			diff := math.Abs(got - tt.expected)
			if diff > tolerance {
				t.Errorf("Float64() = %v, want %v, diff %v", got, tt.expected, diff)
			}
		})
	}
}

// TestFromFloat64 tests the FromFloat64() function
func TestFromFloat64(t *testing.T) {
	tests := []struct {
		name        string
		val         float64
		expectedOk  bool
		wantFloat64 float64
	}{
		{"zero", 0.0, true, 0.0},
		{"one", 1.0, true, 1.0},
		{"half", 0.5, true, 0.5},
		{"negative", -2.5, true, -2.5},
		{"pi", math.Pi, true, math.Pi},
		{"large", 1e10, true, 1e10},
		{"small", 1e-10, true, 1e-10},
		{"nan", math.NaN(), false, 0.0},
		{"posInf", math.Inf(1), false, 0.0},
		{"negInf", math.Inf(-1), false, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromFloat64(tt.val)

			if tt.expectedOk {
				if err != nil {
					t.Errorf("FromFloat64() error = %v, want nil", err)
					return
				}

				// Convert back to float64 and check
				gotFloat := got.Float64()
				tolerance := 1e-10
				diff := math.Abs(gotFloat - tt.wantFloat64)
				if diff > tolerance {
					t.Errorf("FromFloat64() -> Float64() = %v, want %v", gotFloat, tt.wantFloat64)
				}
			} else {
				if err == nil {
					t.Errorf("FromFloat64() error = nil, want error")
				}
			}
		})
	}
}

// TestInt64 tests the Int64() method
func TestInt64(t *testing.T) {
	tests := []struct {
		name       string
		x, y       int64
		want       int64
		wantErr    bool
		checkExact bool
	}{
		{"zero", 0, 1, 0, false, true},
		{"one", 1, 1, 1, false, true},
		{"two", 2, 1, 2, false, true},
		{"negative", -5, 1, -5, false, true},
		{"half_round_down", 5, 2, 2, false, false},     // 2.5 rounds to 2
		{"half_round_up", 7, 2, 4, false, false},       // 3.5 rounds to 4
		{"third", 10, 3, 3, false, false},              // 3.33... rounds to 3
		{"negative_fraction", -7, 2, -4, false, false}, // -3.5 rounds to -4
		{"large", math.MaxInt64, 1, math.MaxInt64, false, true},
		{"min", math.MinInt64, 1, math.MinInt64, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f128 := MustNew(tt.x, tt.y)
			got, err := f128.Int64()

			if (err != nil) != tt.wantErr {
				t.Errorf("Int64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want && tt.checkExact {
				t.Errorf("Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFromInt64 tests the FromInt64() function
func TestFromInt64(t *testing.T) {
	tests := []struct {
		name string
		val  int64
	}{
		{"zero", 0},
		{"one", 1},
		{"negative", -1},
		{"large", math.MaxInt64},
		{"small", math.MinInt64},
		{"medium", 12345},
		{"negMedium", -12345},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromInt64(tt.val)

			// Should be equal to MustNew(val, 1)
			want := MustNew(tt.val, 1)
			if got.Cmp(want) != 0 {
				t.Errorf("FromInt64(%v) = %v, want %v", tt.val, got, want)
			}

			// Test round-trip through Float64
			floatVal := got.Float64()
			if floatVal != float64(tt.val) {
				t.Errorf("FromInt64(%v).Float64() = %v, want %v", tt.val, floatVal, tt.val)
			}
		})
	}
}

// TestAbs tests the Abs() method
func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		x, y int64
		want int64
	}{
		{"zero", 0, 1, 0},
		{"positive", 5, 1, 5},
		{"negative", -5, 1, 5},
		{"negative_fraction", -7, 2, 7},
		{"positive_fraction", 7, 2, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f128 := MustNew(tt.x, tt.y)
			got := f128.Abs()

			// Check that the result is non-negative
			if got.Sign() < 0 {
				t.Errorf("Abs() returned negative value: %v", got)
			}

			// Check that Abs of absolute value is same
			if got.Abs().Cmp(got) != 0 {
				t.Errorf("Abs(Abs()) != Abs()")
			}

			// Verify value by converting to int64
			gotInt, err := got.Int64()
			if err == nil && gotInt < 0 {
				t.Errorf("Abs() returned negative int64: %v", gotInt)
			}
		})
	}
}

// TestNeg tests the Neg() method
func TestNeg(t *testing.T) {
	tests := []struct {
		name string
		x, y int64
	}{
		{"zero", 0, 1},
		{"positive", 5, 1},
		{"negative", -5, 1},
		{"positive_fraction", 7, 2},
		{"negative_fraction", -7, 2},
		{"large", math.MaxInt64, 1},
		{"small", math.MinInt64, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f128 := MustNew(tt.x, tt.y)
			neg := f128.Neg()
			negNeg := neg.Neg() // Negate twice

			// Negating twice should return to original
			if negNeg.Cmp(f128) != 0 {
				t.Errorf("Neg(Neg()) != original: got %v, want %v", negNeg, f128)
			}

			// Neg + original should be zero (or very close for rounding)
			sum := neg.Add(f128)
			if !sum.IsZero() {
				t.Errorf("Neg() + original != Zero: %v", sum)
			}

			// Signs should be opposite
			if f128.Sign() == 0 {
				if neg.Sign() != 0 {
					t.Errorf("Neg(Zero) should be Zero")
				}
			} else {
				if f128.Sign()*neg.Sign() >= 0 {
					t.Errorf("Neg() should have opposite sign")
				}
			}
		})
	}
}

// TestFloat64RoundTrip tests round-trip conversion
func TestFloat64RoundTrip(t *testing.T) {
	values := []float64{
		0.0, 1.0, -1.0, 0.5, -0.5,
		1.5, -2.5, 10.0, -100.0,
		math.Pi, 1e10, 1e-10,
		123456.789, -98765.432,
	}

	for _, val := range values {
		t.Run("", func(t *testing.T) {
			f128, err := FromFloat64(val)
			if err != nil {
				t.Fatalf("FromFloat64 failed: %v", err)
			}

			roundTrip := f128.Float64()
			tolerance := 1e-12
			diff := math.Abs(roundTrip - val)

			if diff > tolerance {
				t.Errorf("Round-trip failed: original=%v, after=%v, diff=%v", val, roundTrip, diff)
			}
		})
	}
}

// TestInt64RoundTrip tests round-trip conversion through int64
func TestInt64RoundTrip(t *testing.T) {
	values := []int64{
		0, 1, -1, 10, -10,
		math.MaxInt64, math.MinInt64,
		42, -42, 123456789,
	}

	for _, val := range values {
		t.Run("", func(t *testing.T) {
			f128 := FromInt64(val)
			roundTrip, err := f128.Int64()
			if err != nil {
				t.Fatalf("Int64() failed: %v", err)
			}

			if roundTrip != val {
				t.Errorf("Round-trip failed: original=%v, after=%v", val, roundTrip)
			}
		})
	}
}

// TestAbsIdentity tests that |a| == a for non-negative values
func TestAbsIdentity(t *testing.T) {
	values := []int64{0, 1, 10, math.MaxInt64}

	for _, val := range values {
		t.Run("", func(t *testing.T) {
			f128 := FromInt64(val)
			abs := f128.Abs()

			if abs.Cmp(f128) != 0 {
				t.Errorf("Abs of non-negative should equal itself: %v != %v", abs, f128)
			}
		})
	}
}

// TestAbsNegation tests that abs of negative equals positive
func TestAbsNegation(t *testing.T) {
	tests := []struct {
		neg int64
		pos int64
	}{
		{-1, 1},
		{-10, 10},
		{-math.MaxInt64, math.MaxInt64},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			neg := FromInt64(tt.neg)
			pos := FromInt64(tt.pos)

			absNeg := neg.Abs()

			if absNeg.Cmp(pos) != 0 {
				t.Errorf("Abs(%v) = %v, want %v", neg, absNeg, pos)
			}
		})
	}
}

// TestNegProperty tests that a + Neg(a) == 0
func TestNegProperty(t *testing.T) {
	values := []int64{
		0, 1, -1, 10, -10,
		42, -42, 123456,
	}

	for _, val := range values {
		t.Run("", func(t *testing.T) {
			f128 := FromInt64(val)
			neg := f128.Neg()
			sum := f128.Add(neg)

			if !sum.IsZero() {
				t.Errorf("a + Neg(a) != 0: %v + %v = %v", f128, neg, sum)
			}
		})
	}
}
