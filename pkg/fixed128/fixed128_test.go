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

func TestAdd(t *testing.T) {
	tt := []struct {
		name string
		a    Fixed128
		b    Fixed128
		want Fixed128
		err  bool
	}{
		// Same signs
		{"positive + positive", Fixed128{hi: 5, lo: 10}, Fixed128{hi: 3, lo: 20}, Fixed128{hi: 8, lo: 30}, false},
		{"negative + negative", Fixed128{hi: 5, lo: 10, neg: true}, Fixed128{hi: 3, lo: 20, neg: true}, Fixed128{hi: 8, lo: 30, neg: true}, false},
		{"zero + positive", Fixed128{}, Fixed128{hi: 5}, Fixed128{hi: 5}, false},
		{"negative + zero", Fixed128{hi: 5, neg: true}, Fixed128{}, Fixed128{hi: 5, neg: true}, false},

		// Different signs
		{"positive + negative (pos larger)", Fixed128{hi: 10}, Fixed128{hi: 3, neg: true}, Fixed128{hi: 7}, false},
		{"negative + positive (neg larger)", Fixed128{hi: 10, neg: true}, Fixed128{hi: 3}, Fixed128{hi: 7, neg: true}, false},
		{"opposite signs, equal magnitude", Fixed128{hi: 5}, Fixed128{hi: 5, neg: true}, Fixed128{}, false},
		{"opposite signs with lo", Fixed128{hi: 5, lo: 100}, Fixed128{hi: 5, lo: 50, neg: true}, Fixed128{lo: 50}, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.a.Add(tc.b)
			if (err != nil) != tc.err {
				t.Errorf("Add() error = %v, want error %v", err, tc.err)
				return
			}
			if got != tc.want {
				t.Errorf("Add(%+v, %+v) = %+v, want %+v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tt := []struct {
		name string
		a    Fixed128
		b    Fixed128
		want Fixed128
		err  bool
	}{
		// Same signs
		{"positive - positive (a larger)", Fixed128{hi: 10}, Fixed128{hi: 3}, Fixed128{hi: 7}, false},
		{"positive - positive (b larger)", Fixed128{hi: 3}, Fixed128{hi: 10}, Fixed128{hi: 7, neg: true}, false},
		{"positive - positive (equal)", Fixed128{hi: 5}, Fixed128{hi: 5}, Fixed128{}, false},
		{"negative - negative (a larger magnitude)", Fixed128{hi: 10, neg: true}, Fixed128{hi: 3, neg: true}, Fixed128{hi: 7, neg: true}, false},
		{"negative - negative (b larger magnitude)", Fixed128{hi: 3, neg: true}, Fixed128{hi: 10, neg: true}, Fixed128{hi: 7}, false},

		// Different signs
		{"positive - negative", Fixed128{hi: 5}, Fixed128{hi: 3, neg: true}, Fixed128{hi: 8}, false},
		{"negative - positive", Fixed128{hi: 5, neg: true}, Fixed128{hi: 3}, Fixed128{hi: 8, neg: true}, false},

		// With zero
		{"positive - zero", Fixed128{hi: 5}, Fixed128{}, Fixed128{hi: 5}, false},
		{"zero - positive", Fixed128{}, Fixed128{hi: 5}, Fixed128{hi: 5, neg: true}, false},
		{"negative - zero", Fixed128{hi: 5, neg: true}, Fixed128{}, Fixed128{hi: 5, neg: true}, false},

		// With lo parts
		{"with lo parts", Fixed128{hi: 5, lo: 100}, Fixed128{hi: 2, lo: 50}, Fixed128{hi: 3, lo: 50}, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.a.Sub(tc.b)
			if (err != nil) != tc.err {
				t.Errorf("Sub() error = %v, want error %v", err, tc.err)
				return
			}
			if got != tc.want {
				t.Errorf("Sub(%+v, %+v) = %+v, want %+v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestAddOverflow(t *testing.T) {
	tt := []struct {
		name string
		a    Fixed128
		b    Fixed128
	}{
		{"max hi + 1", Fixed128{hi: ^uint64(0)}, Fixed128{hi: 1}},
		{"max values", Fixed128{hi: ^uint64(0), lo: ^uint64(0)}, Fixed128{hi: 1}},
		{"positive overflow with lo carry", Fixed128{hi: ^uint64(0), lo: ^uint64(0)}, Fixed128{lo: 1}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.a.Add(tc.b)
			if err == nil {
				t.Errorf("Add(%+v, %+v) expected error, got nil", tc.a, tc.b)
			}
		})
	}
}

func TestSubUnderflow(t *testing.T) {
	tt := []struct {
		name string
		a    Fixed128
		b    Fixed128
	}{
		{"positive small - positive large", Fixed128{hi: 1}, Fixed128{hi: ^uint64(0)}},
		{"same sign, small hi - large hi", Fixed128{hi: 1, neg: true}, Fixed128{hi: ^uint64(0), neg: true}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.a.Sub(tc.b)
			if err == nil {
				t.Errorf("Sub(%+v, %+v) expected error, got nil", tc.a, tc.b)
			}
		})
	}
}

func TestMulInt64(t *testing.T) {
	tt := []struct {
		name       string
		f128       Fixed128
		multiplier int64
		want       int64
		err        bool
	}{
		{"positive * positive", Fixed128{hi: 10}, 5, 50, false},
		{"positive * negative", Fixed128{hi: 10}, -5, -50, false},
		{"negative * positive", Fixed128{hi: 10, neg: true}, 5, -50, false},
		{"negative * negative", Fixed128{hi: 10, neg: true}, -5, 50, false},
		{"zero * any", Fixed128{}, 100, 0, false},
		{"any * zero", Fixed128{hi: 100}, 0, 0, false},
		{"with lo part", Fixed128{hi: 1, lo: 1 << 63}, 2, 3, false},
		{"inverse of division", Fixed128{hi: 1, lo: 0}, 86_400_000_000_000, 86_400_000_000_000, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.f128.MulInt64(tc.multiplier)
			if (err != nil) != tc.err {
				t.Errorf("MulInt64() error = %v, want error %v", err, tc.err)
				return
			}
			if got != tc.want {
				t.Errorf("MulInt64(%+v, %d) = %d, want %d", tc.f128, tc.multiplier, got, tc.want)
			}
		})
	}
}

func TestMulInt64Overflow(t *testing.T) {
	tt := []struct {
		name       string
		f128       Fixed128
		multiplier int64
	}{
		{"max hi * large", Fixed128{hi: ^uint64(0)}, 2},
		{"large * large", Fixed128{hi: 1 << 62}, 1 << 10},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.f128.MulInt64(tc.multiplier)
			if err == nil {
				t.Errorf("MulInt64(%+v, %d) expected error, got nil", tc.f128, tc.multiplier)
			}
		})
	}
}

func TestAbsCmp(t *testing.T) {
	tt := []struct {
		name string
		a    Fixed128
		b    Fixed128
		want int
	}{
		{"equal zero", Fixed128{}, Fixed128{}, 0},
		{"equal non-zero", Fixed128{hi: 5, lo: 10}, Fixed128{hi: 5, lo: 10}, 0},
		{"a > b by hi", Fixed128{hi: 10}, Fixed128{hi: 5}, 1},
		{"a < b by hi", Fixed128{hi: 5}, Fixed128{hi: 10}, -1},
		{"hi equal, a > b by lo", Fixed128{hi: 5, lo: 20}, Fixed128{hi: 5, lo: 10}, 1},
		{"hi equal, a < b by lo", Fixed128{hi: 5, lo: 10}, Fixed128{hi: 5, lo: 20}, -1},
		{"hi >, lo<", Fixed128{hi: 10, lo: 5}, Fixed128{hi: 5, lo: 10}, 1},
		{"hi <, lo>", Fixed128{hi: 5, lo: 10}, Fixed128{hi: 10, lo: 5}, -1},
		{"large values equal", Fixed128{hi: 1<<63 - 1, lo: 1<<63 - 1}, Fixed128{hi: 1<<63 - 1, lo: 1<<63 - 1}, 0},
		{"large hi diff", Fixed128{hi: 1 << 63}, Fixed128{hi: 1}, 1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := absCmp(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("absCmp(%+v, %+v) = %d, want %d", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
