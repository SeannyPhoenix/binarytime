package binarytime

import "testing"

func BenchmarkNow(b *testing.B) {
	for b.Loop() {
		_ = Now()
	}
}
