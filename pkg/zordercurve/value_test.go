package zordercurve

import (
	"testing"
)

func FuzzXYToValToXY(f *testing.F) {
	tt := []struct {
		x uint32
		y uint32
	}{
		{3, 4},
		{7, 5},
		{12345, 987654321},
		{0xffffffff, 34},
		{656, 0xffffffff},
	}

	for _, tc := range tt {
		f.Add(tc.x, tc.y)
	}

	f.Fuzz(func(t *testing.T, x, y uint32) {
		v := GetValueFromXY(x, y)
		cx, cy := GetXYFromValue(v)
		if cx != x || cy != y {
			t.Fatalf("result %d, %d does not equal given %d, %d", cx, cy, x, y)
		}
	})
}

func FuzzValToXYToVal(f *testing.F) {
	tt := []uint64{1, 3625673, 0xFFFFFFFFFFFFFFFF, 0xabcdef0123456789}

	for _, v := range tt {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, v uint64) {
		x, y := GetXYFromValue(v)
		cv := GetValueFromXY(x, y)
		if cv != v {
			t.Fatalf("result %d does not equal given %d", cv, v)
		}
	})
}
