package binarytime

import "testing"

func BenchmarkNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Now()
	}
}
