package genetic

import "testing"

func BenchmarkRandom(b *testing.B) {
	b.Run("randInt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			randInt(1000)
		}
	})
	b.Run("randFloat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			randFloat()
		}
	})
}
