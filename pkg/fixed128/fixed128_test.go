package fixed128

import (
	"testing"
)

func FuzzFixed128(t *testing.F) {
	tt := []struct {
		dividend int64
		divisor  uint64
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

	t.Fuzz(func(t *testing.T, dividend int64, divisor uint64) {
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
