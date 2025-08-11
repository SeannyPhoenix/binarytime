package fixed128

import (
	"fmt"
	"testing"
)

func TestFized128(t *testing.T) {
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

	for _, test := range tt {
		t.Run(fmt.Sprintf("%d / %d", test.dividend, test.divisor), func(t *testing.T) {
			f128, err := NewF128(test.dividend, test.divisor)
			if err != nil {
				t.Fatalf("failed to create Fixed128: %v", err)
			}

			got, err := f128.FromF128(test.divisor)
			if err != nil {
				t.Fatalf("failed to convert from Fixed128: %v", err)
			}

			if got != test.dividend {
				t.Errorf("unexpected result: got %d, want %d", got, test.dividend)
			}
		})
	}
}
