package fixed128

import (
	"testing"
)

func FuzzFixed128(t *testing.F) {
	tt := []struct {
		dividend int64
		divisor  int64
	}{{
		dividend: 12345,
		divisor:  37,
	}, {
		dividend: -36756375637985,
		divisor:  2437523678,
	}}

	for _, tc := range tt {
		t.Add(tc.dividend, tc.divisor)
	}

	t.Fuzz(func(t *testing.T, dividend int64, divisor int64) {
		if divisor == 0 {
			t.Skip("skipping test case with division by zero")
		}

		f128, err := NewF128(dividend, divisor)
		if err != nil {
			t.Fatalf("failed to create Fixed128: %v", err)
		}

		got, err := f128.FromF128(divisor)
		if err != nil {
			t.Fatalf("failed to convert from Fixed128: %v", err)
		}

		if got != dividend {
			t.Errorf("unexpected result: got %d, want %d", got, dividend)
		}
	})
}

func BenchmarkNewF128(b *testing.B) {
	for b.Loop() {
		_, err := NewF128(1234567890123456789, 987654321)
		if err != nil {
			b.Fatalf("failed to create Fixed128: %v", err)
		}
	}
}

func BenchmarkFrom128(b *testing.B) {
	f128, err := NewF128(1234567890123456789, 987654321)
	if err != nil {
		b.Fatalf("failed to create Fixed128: %v", err)
	}

	for b.Loop() {
		_, err := f128.FromF128(987654321)
		if err != nil {
			b.Fatalf("failed to convert from Fixed128: %v", err)
		}
	}
}
