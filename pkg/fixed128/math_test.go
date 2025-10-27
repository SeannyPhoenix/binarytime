package fixed128

import (
	"fmt"
	"math"
	"testing"
)

// TestMulBasic tests basic multiplication operations
func TestMulBasic(t *testing.T) {
	tests := []struct {
		name string
		a    Fixed128
		b    Fixed128
		want Fixed128
	}{
		{"2*3", MustNew(2, 1), MustNew(3, 1), MustNew(6, 1)},
		{"3*4", MustNew(3, 1), MustNew(4, 1), MustNew(12, 1)},
		{"1*1", MustNew(1, 1), MustNew(1, 1), MustNew(1, 1)},
		{"0*5", MustNew(0, 1), MustNew(5, 1), MustNew(0, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Mul(tt.b)
			if got.Cmp(tt.want) != 0 {
				t.Errorf("Mul(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TestMulFractional tests fractional multiplication
func TestMulFractional(t *testing.T) {
	tests := []struct {
		name string
		a    Fixed128
		b    Fixed128
		want Fixed128
	}{
		{
			name: "1/2 * 1/2 = 1/4",
			a:    MustNew(1, 2),
			b:    MustNew(1, 2),
			want: MustNew(1, 4),
		},
		{
			name: "1/3 * 1/3 = 1/9",
			a:    MustNew(1, 3),
			b:    MustNew(1, 3),
			want: MustNew(1, 9),
		},
		{
			name: "1/8 * 1/2 = 1/16",
			a:    MustNew(1, 8),
			b:    MustNew(1, 2),
			want: MustNew(1, 16),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Mul(tt.b)
			if got.Cmp(tt.want) != 0 {
				t.Errorf("Mul(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TestQuoBasic tests basic division operations
func TestQuoBasic(t *testing.T) {
	tests := []struct {
		name    string
		a       int64
		b       int64
		c       int64
		d       int64
		want    int64
		wantErr bool
	}{
		{"6/2", 6, 1, 2, 1, 3, false},
		{"12/3", 12, 1, 3, 1, 4, false},
		{"1/1", 1, 1, 1, 1, 1, false},
		{"0/5", 0, 1, 5, 1, 0, false},
		{"division by zero", 5, 1, 0, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := MustNew(tt.a, 1)
			if tt.b != 1 {
				a = MustNew(tt.a, tt.b)
			}

			var b Fixed128
			if tt.c == 0 && tt.d == 1 {
				// Test division by zero
				b = Zero
			} else {
				b = MustNew(tt.c, 1)
				if tt.d != 1 {
					b = MustNew(tt.c, tt.d)
				}
			}

			got, err := a.Quo(b)

			if (err != nil) != tt.wantErr {
				t.Errorf("Quo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				want := MustNew(tt.want, 1)
				if got.Cmp(want) != 0 {
					t.Errorf("Quo(%v, %v) = %v, want %v", a, b, got, want)
				}
			}
		})
	}
}

// TestQuoFractional tests fractional division
func TestQuoFractional(t *testing.T) {
	tests := []struct {
		name string
		a    int64
		b    int64
		c    int64
		d    int64
		want int64
	}{
		{"1/2 / 1/4", 1, 2, 1, 4, 2},
		{"1/2 / 1/8", 1, 2, 1, 8, 4},
		{"3/2 / 1/4", 3, 2, 1, 4, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := MustNew(tt.a, tt.b)
			b := MustNew(tt.c, tt.d)
			want := MustNew(tt.want, 1)

			got, err := a.Quo(b)
			if err != nil {
				t.Fatalf("Quo() error = %v", err)
			}

			if got.Cmp(want) != 0 {
				t.Errorf("Quo(%v, %v) = %v, want %v", a, b, got, want)
			}
		})
	}
}

// TestMulIdentity tests multiplication identity property
func TestMulIdentity(t *testing.T) {
	values := []int64{1, 2, 3, 10, 100, -1, -5}
	for _, v := range values {
		t.Run(fmt.Sprintf("1*%d", v), func(t *testing.T) {
			val := MustNew(v, 1)
			one := MustNew(1, 1)
			got := one.Mul(val)
			if got.Cmp(val) != 0 {
				t.Errorf("1 * %v = %v, want %v", val, got, val)
			}
		})
	}
}

// TestQuoIdentity tests division identity property
func TestQuoIdentity(t *testing.T) {
	values := []int64{1, 2, 3, 10, 100, -1, -5}
	for _, v := range values {
		t.Run(fmt.Sprintf("%d/1", v), func(t *testing.T) {
			val := MustNew(v, 1)
			one := MustNew(1, 1)
			got, err := val.Quo(one)
			if err != nil {
				t.Fatalf("Quo() error = %v", err)
			}
			if got.Cmp(val) != 0 {
				t.Errorf("%v / 1 = %v, want %v", val, got, val)
			}
		})
	}
}

// TestMulDistributive tests distributive property: a * (b + c) = a*b + a*c
func TestMulDistributive(t *testing.T) {
	a := MustNew(3, 4) // 0.75
	b := MustNew(1, 4) // 0.25
	c := MustNew(1, 8) // 0.125

	left := a.Mul(b.Add(c))
	right := a.Mul(b).Add(a.Mul(c))

	if left.Cmp(right) != 0 {
		t.Errorf("a * (b + c) = %v, want %v", left, right)
	}
}

// TestLargeNumbers tests with large numbers
func TestLargeNumbers(t *testing.T) {
	// Test near max int64
	a := MustNew(math.MaxInt64, 1)
	b := MustNew(2, 1)

	// This should work without overflow
	got := a.Mul(b)
	if got.IsZero() {
		t.Fatal("Mul with large numbers resulted in zero")
	}

	// Test division by large numbers
	c := MustNew(math.MaxInt64, 2)
	_, err := a.Quo(c)
	if err != nil {
		t.Fatalf("Quo() error = %v", err)
	}
}
