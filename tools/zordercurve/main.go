package main

import (
	"fmt"

	"github.com/seannyphoenix/binarytime/pkg/zordercurve"
)

func main() {
	toRun := []string{
		// "a",
		"b",
	}

	for _, k := range toRun {
		if f, ok := funcs[k]; ok {
			f()
		}
	}
}

var funcs = map[string]func(){
	"a": func() {
		fmt.Println(zordercurve.GetValueFromXY(7, 5))
		fmt.Println(zordercurve.GetXYFromValue(0x99))
	},

	"b": func() {
		for x := uint32(0); x < 16; x++ {
			c := uint32(1 << (2 * x))
			v := zordercurve.GetValueFromXY(c, 0)
			s := uint64(c) * uint64(c)
			fmt.Printf("%02d, %10d: %20d = %20d %064b\n", x, c, v, s, v)
		}
	},
}
