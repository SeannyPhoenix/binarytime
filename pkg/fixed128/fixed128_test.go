package fixed128

import (
	"fmt"
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

func BenchmarkNew(b *testing.B) {
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

// Test negate overflow on math.MinInt64 during FromF128
func TestFromF128_NegateOverflow(t *testing.T) {
	// divisor = 1 gives hi = x, lo = 0
	f128, err := New(math.MinInt64, 1)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// FromF128 with negative divisor flips sign,
	// which would try to negate math.MinInt64
	_, err = f128.MulInt64(-1)
	if err == nil {
		t.Fatalf("expected overflow error when negating math.MinInt64")
	}
}

// Make sure << (64 - i) isn't undefined
func TestGetComponents_ShiftSafety(t *testing.T) {
	_, lo := getComponents(1, 2) // 0.5 => lo should be >= 1<<63
	if lo == 0 {
		t.Fatalf("unexpected zero lo for x=1,y=2")
	}
}

// Basic round-trip of MarshalBinary/UnmarshalBinary
func TestMarshalBinaryRoundTrip(t *testing.T) {
	vals := []Fixed128{
		MustNew(123, 7),
		MustNew(-123, 7),
		MustNew(0, 1),
		MustNew(math.MaxInt64, 3),
	}

	for _, v := range vals {
		bin, err := v.MarshalBinary()
		if err != nil {
			t.Errorf("marshal failed: %v", err)
			continue
		}

		var out Fixed128
		if err := out.UnmarshalBinary(bin); err != nil {
			t.Errorf("unmarshal failed: %v", err)
			continue
		}

		if out.Cmp(v) != 0 {
			t.Errorf("round-trip mismatch: want %v got %v", v, out)
		}
	}
}

// Test hydrate rounding behaviour on divisors without power-of-two scale
func TestHydrateRounding(t *testing.T) {
	type tc struct {
		lo, div uint64
		want    uint64
	}

	// Hand-crafted small cases
	tests := []tc{
		{lo: 1<<63 - 1, div: 3, want: 1}, // .9999... * 1/3 ≈ 0.33 => round
		{lo: 1 << 63, div: 3, want: 1},
		{lo: 0, div: 5, want: 0},
	}

	for _, c := range tests {
		t.Run(fmt.Sprintf("lo=%d,div=%d", c.lo, c.div), func(t *testing.T) {
			got := hydrate(c.lo, c.div)
			if got != c.want {
				t.Errorf("hydrate(%d,%d) = %d, want %d", c.lo, c.div, got, c.want)
			}
		})
	}
}

// MarshalText -> UnmarshalText -> round-trip
func TestTextRoundTrip(t *testing.T) {
	tests := []Fixed128{
		MustNew(10, 3),
		MustNew(-10, 3),
		MustNew(123456789, 987654321),
	}

	for _, v := range tests {
		b, err := v.MarshalText()
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}
		var out Fixed128
		if err := out.UnmarshalText(b); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if v.Cmp(out) != 0 {
			t.Fatalf("round-trip mismatch: want=%v got=%v", v, out)
		}
	}
}
