package fixed128

import (
	"math"
	"testing"
)

// TestDecimalString tests the DecimalString() method
func TestDecimalString(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int64
		expected string
	}{
		{"zero", 0, 1, "0.000000000000000"},
		{"one", 1, 1, "1.000000000000000"},
		{"half", 1, 2, "0.500000000000000"},
		{"third", 1, 3, "0.333333333333333"},
		{"quarter", 1, 4, "0.250000000000000"},
		{"pi_approx", 22, 7, "3.142857142857143"},
		{"negative", -5, 2, "-2.500000000000000"},
		{"large", 1000000, 1, "1000000.000000000000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f128 := MustNew(tt.x, tt.y)
			got := f128.DecimalString()
			if got != tt.expected {
				t.Errorf("DecimalString() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestDecimalStringWithPrecision tests custom precision
func TestDecimalStringWithPrecision(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int64
		prec     int
		expected string
	}{
		{"zero0", 0, 1, 0, "0"},
		{"zero5", 0, 1, 5, "0.00000"},
		{"one0", 1, 1, 0, "1"},
		{"one2", 1, 1, 2, "1.00"},
		{"half3", 1, 2, 3, "0.500"},
		{"third1", 1, 3, 1, "0.3"},
		{"third5", 1, 3, 5, "0.33333"},
		{"pi0", 22, 7, 0, "3"},
		{"pi5", 22, 7, 5, "3.14286"},
		{"neg0", -5, 2, 0, "-2"},
		{"neg2", -5, 2, 2, "-2.50"},
		{"negative", -10, 3, 1, "-3.3"},
		{"large", 123456, 1, 0, "123456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f128 := MustNew(tt.x, tt.y)
			got := f128.DecimalStringWithPrecision(tt.prec)
			if got != tt.expected {
				t.Errorf("DecimalStringWithPrecision(%d) = %q, want %q", tt.prec, got, tt.expected)
			}
		})
	}
}

// TestDecimalStringExtremeValues tests edge cases
func TestDecimalStringExtremeValues(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int64
		validate func(string) bool
	}{
		{
			name: "maxint",
			x:    math.MaxInt64,
			y:    1,
			validate: func(s string) bool {
				// Should start with 9.22337...
				return s[0] == '9' && len(s) > 2
			},
		},
		{
			name: "minint",
			x:    math.MinInt64,
			y:    1,
			validate: func(s string) bool {
				// Should start with -9.22337...
				return s[0] == '-' && s[1] == '9'
			},
		},
		{
			name: "very_small",
			x:    1,
			y:    1000000,
			validate: func(s string) bool {
				return s[:2] == "0."
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f128 := MustNew(tt.x, tt.y)
			got := f128.DecimalString()
			if !tt.validate(got) {
				t.Errorf("DecimalString() = %q, validation failed", got)
			}
		})
	}
}

// TestDecimalStringPrecision0 tests that precision 0 returns integer
func TestDecimalStringPrecision0(t *testing.T) {
	tests := []struct {
		x, y     int64
		expected string
	}{
		{0, 1, "0"},
		{5, 1, "5"},
		{-5, 1, "-5"},
		{7, 2, "4"},   // 3.5 rounds to 4
		{5, 2, "2"},   // 2.5 rounds to 2
		{-7, 2, "-4"}, // -3.5 rounds to -4
		{-5, 2, "-2"}, // -2.5 rounds to -2
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f128 := MustNew(tt.x, tt.y)
			got := f128.DecimalStringWithPrecision(0)
			if got != tt.expected {
				t.Errorf("DecimalStringWithPrecision(0) = %q, want %q", got, tt.expected)
			}
		})
	}
}
