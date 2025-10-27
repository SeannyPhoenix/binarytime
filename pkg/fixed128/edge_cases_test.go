package fixed128

import (
	"math"
	"testing"
)

// TestAddEdgeCases tests edge cases for addition
func TestAddEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		a, b Fixed128
		want Fixed128
	}{
		{"zero + zero", Zero, Zero, Zero},
		{"zero + positive", Zero, MustNew(5, 1), MustNew(5, 1)},
		{"positive + zero", MustNew(5, 1), Zero, MustNew(5, 1)},
		{"negative + positive equals zero", MustNew(-5, 1), MustNew(5, 1), Zero},
		{"positive + negative equals zero", MustNew(5, 1), MustNew(-5, 1), Zero},
		{"maxint + zero", MustNew(math.MaxInt64, 1), Zero, MustNew(math.MaxInt64, 1)},
		{"minint + zero", MustNew(math.MinInt64, 1), Zero, MustNew(math.MinInt64, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Add(tt.b)
			if got.Cmp(tt.want) != 0 {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSubEdgeCases tests edge cases for subtraction
func TestSubEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		a, b Fixed128
		want Fixed128
	}{
		{"zero - zero", Zero, Zero, Zero},
		{"positive - zero", MustNew(5, 1), Zero, MustNew(5, 1)},
		{"positive - same", MustNew(5, 1), MustNew(5, 1), Zero},
		{"positive - larger", MustNew(5, 1), MustNew(10, 1), MustNew(-5, 1)},
		{"negative - zero", MustNew(-5, 1), Zero, MustNew(-5, 1)},
		{"zero - positive", Zero, MustNew(5, 1), MustNew(-5, 1)},
		{"minint - zero", MustNew(math.MinInt64, 1), Zero, MustNew(math.MinInt64, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Sub(tt.b)
			if got.Cmp(tt.want) != 0 {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMulEdgeCases tests edge cases for multiplication
func TestMulEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		a, b Fixed128
		want Fixed128
	}{
		{"zero * zero", Zero, Zero, Zero},
		{"zero * positive", Zero, MustNew(5, 1), Zero},
		{"positive * zero", MustNew(5, 1), Zero, Zero},
		{"one * positive", One, MustNew(5, 1), MustNew(5, 1)},
		{"positive * one", MustNew(5, 1), One, MustNew(5, 1)},
		{"negative * positive", MustNew(-5, 1), MustNew(2, 1), MustNew(-10, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Mul(tt.b)
			if got.Cmp(tt.want) != 0 {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestQuoEdgeCases tests edge cases for division
func TestQuoEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		a, b    Fixed128
		want    Fixed128
		wantErr bool
	}{
		{"zero / positive", Zero, MustNew(5, 1), Zero, false},
		{"positive / one", MustNew(5, 1), One, MustNew(5, 1), false},
		{"positive / zero", MustNew(5, 1), Zero, Fixed128{}, true},
		{"negative / positive", MustNew(-10, 1), MustNew(2, 1), MustNew(-5, 1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Quo(tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Quo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Cmp(tt.want) != 0 {
				t.Errorf("Quo() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestArithmeticProperties tests mathematical properties
func TestArithmeticProperties(t *testing.T) {
	t.Run("Identity property: a + 0 = a", func(t *testing.T) {
		values := []Fixed128{
			MustNew(1, 1), MustNew(5, 1), MustNew(-3, 1),
			MustNew(1, 2), MustNew(22, 7),
		}
		for _, val := range values {
			got := val.Add(Zero)
			if got.Cmp(val) != 0 {
				t.Errorf("a + 0 = %v, want %v", got, val)
			}
		}
	})

	t.Run("Additive inverse: a + (-a) = 0", func(t *testing.T) {
		values := []Fixed128{
			MustNew(1, 1), MustNew(5, 1), MustNew(-3, 1),
			MustNew(1, 2), MustNew(22, 7),
		}
		for _, val := range values {
			got := val.Add(val.Neg())
			if !got.IsZero() {
				t.Errorf("a + (-a) = %v, want 0", got)
			}
		}
	})

	t.Run("Associative property: (a + b) - b = a", func(t *testing.T) {
		a := MustNew(22, 7)
		b := MustNew(1, 3)
		got := a.Add(b).Sub(b)
		if got.Cmp(a) != 0 {
			t.Errorf("(a + b) - b = %v, want %v", got, a)
		}
	})

	t.Run("Distributive property: a * (b + c) = a*b + a*c", func(t *testing.T) {
		a := MustNew(1, 2)
		b := MustNew(1, 3)
		c := MustNew(1, 4)

		left := a.Mul(b.Add(c))
		right := a.Mul(b).Add(a.Mul(c))

		if left.Cmp(right) != 0 {
			t.Errorf("a * (b + c) = %v, want %v", left, right)
		}
	})
}

// TestSignEdgeCases tests sign handling with edge cases
func TestSignEdgeCases(t *testing.T) {
	t.Run("Zero has no sign", func(t *testing.T) {
		if Zero.Sign() != 0 {
			t.Errorf("Zero.Sign() = %d, want 0", Zero.Sign())
		}
		if !Zero.IsZero() {
			t.Errorf("Zero.IsZero() = false, want true")
		}
		if Zero.IsNeg() {
			t.Errorf("Zero.IsNeg() = true, want false")
		}
	})

	t.Run("Positive values", func(t *testing.T) {
		pos := MustNew(5, 1)
		if pos.Sign() <= 0 {
			t.Errorf("pos.Sign() = %d, want > 0", pos.Sign())
		}
		if pos.IsNeg() {
			t.Errorf("pos.IsNeg() = true, want false")
		}
	})

	t.Run("Negative values", func(t *testing.T) {
		neg := MustNew(-5, 1)
		if neg.Sign() >= 0 {
			t.Errorf("neg.Sign() = %d, want < 0", neg.Sign())
		}
		if !neg.IsNeg() {
			t.Errorf("neg.IsNeg() = false, want true")
		}
	})
}

// TestComparisonEdgeCases tests comparison with edge cases
func TestComparisonEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		a, b Fixed128
		want int
	}{
		{"zero = zero", Zero, Zero, 0},
		{"positive > zero", MustNew(5, 1), Zero, 1},
		{"zero < positive", Zero, MustNew(5, 1), -1},
		{"positive < negative", MustNew(5, 1), MustNew(-5, 1), 1},
		{"negative < positive", MustNew(-5, 1), MustNew(5, 1), -1},
		{"equal fractions", MustNew(1, 2), MustNew(2, 4), 0},
		{"small difference", MustNew(1, 1000000), Zero, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Cmp(tt.b)
			if got != tt.want {
				t.Errorf("Cmp() = %d, want %d", got, tt.want)
			}
		})
	}
}

// TestPrecisionLoss tests scenarios where precision might be lost
func TestPrecisionLoss(t *testing.T) {
	t.Run("Very small fractions", func(t *testing.T) {
		// These should work without precision loss
		small := MustNew(1, 1000000)
		expected := MustNew(1, 1000000)

		// Should be exactly equal
		if small.Cmp(expected) != 0 {
			t.Errorf("Small fraction lost precision: got %v, want %v", small, expected)
		}
	})

	t.Run("Large whole numbers", func(t *testing.T) {
		// MaxInt64 should work
		large := MustNew(math.MaxInt64, 1)
		if large.IsZero() {
			t.Errorf("Large number became zero")
		}

		// Should be able to retrieve back
		intVal, err := large.Int64()
		if err != nil {
			t.Errorf("Failed to get Int64(): %v", err)
		} else if intVal != math.MaxInt64 {
			t.Errorf("Retrieved int64 = %d, want %d", intVal, math.MaxInt64)
		}
	})
}
