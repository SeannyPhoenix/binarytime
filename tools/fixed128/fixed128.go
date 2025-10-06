package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/seannyphoenix/binarytime/pkg/binarytime"
	"github.com/seannyphoenix/binarytime/pkg/fixed128"
)

func main() {
	// funcs["a"]()
	// funcs["b"]()
	funcs["c"]()
	funcs["d"]()
}

var funcs = map[string]func(){
	"a": a,
	"b": b,
	"c": c,
	"d": d,
}

func d() {
	f128 := fixed128.MustNew(123478392543, 134332)
	// f128 := fixed128.MustNew(-30, 12)
	b, err := json.MarshalIndent(f128, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println("JSON output:")
	fmt.Println(string(b))
	fmt.Println("String output:")
	fmt.Println(f128)
}

func c() {
	v1 := fixed128.MustNew(30, 6)
	fmt.Println("Fixed128 value:", v1)
}

func b() {
	nbt := binarytime.Now()
	fmt.Println("Current binary time:", nbt)
	fmt.Println("Current binary time:", nbt.Fixed128())
	fmt.Println("Current binary time:", nbt.Fixed128().Value())
}

func a() {
	x, y, err := getArguments()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Arguments received:", x, y)
	f128, err := fixed128.New(x, y)
	if err != nil {
		fmt.Println("Error creating Fixed128:", err)
		return
	}
	fmt.Printf("Fixed128 representation: %s\n", f128)

	r, err := f128.MulInt64(y)
	if err != nil {
		fmt.Println("Error converting from Fixed128:", err)
		return
	}
	fmt.Println("Converted back:", r)
}

func getArguments() (int64, int64, error) {
	var x, y int64
	var err error

	if len(os.Args) < 3 {
		return 0, 0, fmt.Errorf("not enough arguments")
	}

	x, err = strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid first argument: %v", err)
	}

	y, err = strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid second argument: %v", err)
	}

	return x, y, nil
}
