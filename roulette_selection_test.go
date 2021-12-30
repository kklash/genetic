package genetic

import (
	"testing"
)

func TestRouletteSpin(t *testing.T) {
	ticketCounts := []int{
		5000,
		1000,
		100,
		30,
		5,
		0,
	}

	winCounts := make([]int, len(ticketCounts))
	proportions := computeProportions(ticketCounts)
	for i := 0; i < 10000; i++ {
		winner := rouletteSpin(proportions)
		winCounts[winner] += 1
	}

	for i := 0; i < len(ticketCounts)-1; i++ {
		if winCounts[i] <= winCounts[i+1] {
			t.Errorf("expected entrant %d to win more than %d", i, i+1)
		}
	}
}

func BenchmarkRouletteSpin(b *testing.B) {
	ticketCounts := make([]int, 100)
	for i := 0; i < len(ticketCounts); i++ {
		ticketCounts[i] = randInt(50000)
	}

	proportions := computeProportions(ticketCounts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rouletteSpin(proportions)
	}
}
