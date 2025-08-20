package main

import (
	"testing"
)

func BenchmarkRegIsNeg(b *testing.B) {
	var num int64 = -2847

	for i := 0; i < b.N; i++ {
		regIsNeg(num)
	}
}

func BenchmarkBitIsNeg(b *testing.B) {
	var num int64 = -2847

	for i := 0; i < b.N; i++ {
		bitIsNeg(num)
	}
}
