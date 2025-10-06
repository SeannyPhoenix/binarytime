package main

import (
	"fmt"
	"time"

	"github.com/seannyphoenix/binarytime/pkg/binarytime"
)

func main() {
	funcs["a"]()
}

var funcs = map[string]func(){
	"a": a,
}

func a() {
	fmt.Println(binarytime.FromDuration(8 * time.Hour).String())
}
