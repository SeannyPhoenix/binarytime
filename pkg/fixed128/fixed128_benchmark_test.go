package fixed128

import (
	"testing"
)

func BenchmarkCmp(b *testing.B) {
	a := Fixed128{hi: 100, lo: 200, neg: false}
	c := Fixed128{hi: 50, lo: 100, neg: true}
	for b.Loop() {
		_ = a.Cmp(c)
	}
}

func BenchmarkAdd(b *testing.B) {
	a := Fixed128{hi: 100, lo: 200, neg: false}
	c := Fixed128{hi: 50, lo: 100, neg: false}
	for b.Loop() {
		_, _ = a.Add(c)
	}
}

func BenchmarkAddDifferentSigns(b *testing.B) {
	a := Fixed128{hi: 100, lo: 200, neg: false}
	c := Fixed128{hi: 50, lo: 100, neg: true}
	for b.Loop() {
		_, _ = a.Add(c)
	}
}

func BenchmarkSub(b *testing.B) {
	a := Fixed128{hi: 100, lo: 200, neg: false}
	c := Fixed128{hi: 50, lo: 100, neg: false}
	for b.Loop() {
		_, _ = a.Sub(c)
	}
}
