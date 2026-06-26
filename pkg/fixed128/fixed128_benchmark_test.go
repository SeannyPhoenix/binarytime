package fixed128

import (
	"math/rand"
	"testing"
)

func BenchmarkByDivision(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		x, y := rand.Int63(), rand.Int63()
		b.StartTimer()
		_, err := ByDivision(x, y)
		if err != nil {
			b.Fatal(err)
		}
	}
}
