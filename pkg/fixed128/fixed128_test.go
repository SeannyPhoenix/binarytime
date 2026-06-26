package fixed128

import "testing"

func TestFixed128_FromParts(t *testing.T) {
	tt := []struct {
		hi  uint64
		lo  uint64
		neg bool
	}{
		{0, 0, false},
		{1, 1, true},
		{1 << 63, 1 << 63, false},
		{1<<63 - 1, 1<<63 - 1, true},
	}

	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			h, l, n := FromParts(tc.hi, tc.lo, tc.neg).Parts()
			if h != tc.hi || l != tc.lo || n != tc.neg {
				t.Errorf("FromParts(%d, %d, %t) = (%d, %d, %t), want (%d, %d, %t)",
					tc.hi, tc.lo, tc.neg,
					h, l, n,
					tc.hi, tc.lo, tc.neg)
			}
		})
	}
}

func TestFixed128_ByDivision(t *testing.T) {
	tt := []struct {
		x, y int64
		neg  bool
	}{
		{1, 1, false},
		{10, 2, false},
		{-10, 2, true},
		{-10, -2, false},
	}

	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			f128, err := ByDivision(tc.x, tc.y)
			if err != nil {
				t.Fatalf("ByDivision(%d, %d) returned error: %v", tc.x, tc.y, err)
			}
			h, l, n := f128.Parts()
			if n != tc.neg {
				t.Errorf("ByDivision(%d, %d) = (%d, %d, %t), got neg=%t but expected neg=%t",
					tc.x, tc.y,
					h, l, n,
					n, tc.neg)
			}
			if n && h == 0 && l == 0 {
				t.Errorf("ByDivision(%d, %d) = (%d, %d, %t), got zero but expected non-zero",
					tc.x, tc.y,
					h, l, n)
			}
		})
	}
}

func FuzzFixed128_FromParts(f *testing.F) {
	f.Add(uint64(0), uint64(0), false)
	f.Add(uint64(1), uint64(1), true)
	f.Add(uint64(1<<63), uint64(1<<63), false)
	f.Add(uint64(1<<63-1), uint64(1<<63-1), true)

	f.Fuzz(func(t *testing.T, hi, lo uint64, neg bool) {
		h, l, n := FromParts(hi, lo, neg).Parts()
		if h != hi || l != lo || n != neg {
			t.Errorf("FromParts(%d, %d, %t) = (%d, %d, %t), want (%d, %d, %t)",
				hi, lo, neg,
				h, l, n,
				hi, lo, neg)
		}
	})
}
