package main

import "fmt"

func main() {
	funcs["b"]()
}

var funcs = map[string]func(){
	"a": a,
	"b": b,
}

func b() {
	fmt.Printf("%02x\n", 2)
	fmt.Printf("%0[2]*[1]x\n", 259, 4)
}

func a() {
	b := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	fmt.Println(string(b[1:9]))
	fmt.Println(string(b[9:17]))

	b2 := []byte{0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0}

	high := 1
	for {
		byte := b2[high]
		if byte != 0 && high < 9 {
			break
		}
		high++
	}

	low := 17
	for {
		byte := b2[low-1]
		if byte != 0 && low-1 > 9 {
			break
		}
		low--
	}

	fmt.Println(string(b2[high:9]))
	fmt.Println(string(b2[9:low]))
}
