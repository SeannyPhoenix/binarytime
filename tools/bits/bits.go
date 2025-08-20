package main

import (
	"fmt"
	"math"
)

func main() {
	d()
}

func d() {
	a := uint8(0b11110000)
	b := int8(a)
	c := int16(b)
	d := uint16(c)

	fmt.Printf("a: %08b\nb: %08b\nc: %016b\nd: %016b\n", a, b, c, d)
}

func c() {
	i := int16(-0xa5)
	ui := uint16(i)
	fmt.Printf("i  : %016b\nui : %016b\n", i, ui)

	ui = uint16(0xc472)
	i = int16(ui)
	fmt.Printf("ui : %016b\ni  : %016b\n", ui, i)
}

func b() {
	// v := int64(math.MaxInt64)
	// v := int64(math.MinInt64)
	// v := int64(12345)
	v := int64(-1024 - 128 - 16)
	f := uint64(v)
	mask := uint64(v >> 63)
	abs := (uint64(v) ^ mask) - mask
	neg := mask != 0
	fmt.Printf("v  : %064b\nf   : %064b\nmask: %064b\nmag : %064b\nneg : %v\n", v, f, mask, abs, neg)
}

func a() {
	x := int8(-127)
	y := uint8(x)
	z := uint8(-x)
	q := -x
	u := ^x + 1

	fmt.Printf("x: %08b\ny: %08b\nz: %08b\nq: %08b\nu: %08b\n", x, y, z, q, u)
	fmt.Printf("MaxInt8: %08b\nMinInt8 + 1: %08b\n", math.MaxInt8+1, math.MinInt8)
}

func bitIsNeg(v int64) bool {
	return v>>63&1 != 0
}

func regIsNeg(v int64) bool {
	return v < 0
}
