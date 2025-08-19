package fixed128

import (
	"math"
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

		f128, err := New(dividend, divisor)
		if err != nil {
			t.Fatalf("failed to create Fixed128: %v", err)
		}

		got, err := f128.MulInt64(divisor)
		if err != nil {
			t.Fatalf("failed to convert from Fixed128: %v", err)
		}

		if got != dividend {
			t.Errorf("unexpected result: dividend %d, divisor %d, got %d", dividend, divisor, got)
		}
	})
}

func BenchmarkNewF128(b *testing.B) {
	for b.Loop() {
		_, err := New(int64(math.MaxInt64-b.N), int64(b.N+1))
		if err != nil {
			b.Fatalf("failed to create Fixed128: %v", err)
		}
	}
}

func BenchmarkFrom128(b *testing.B) {
	f128, err := New(int64(b.N), int64(math.MaxInt64-b.N))
	if err != nil {
		b.Fatalf("failed to create Fixed128: %v", err)
	}

	for b.Loop() {
		_, err := f128.MulInt64(int64(b.N + 1))
		if err != nil {
			b.Fatalf("failed to convert from Fixed128: %v", err)
		}
	}
}
