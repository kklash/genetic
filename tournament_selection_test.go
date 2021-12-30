package genetic

import (
	"fmt"
	"testing"
)

func makeRandomFitnesses(n int) []int {
	fitnesses := make([]int, n)
	for i := 0; i < len(fitnesses); i++ {
		fitnesses[i] = randInt(50000)
	}
	return fitnesses
}

func benchRandomTournamentWinner(b *testing.B, populationSize, poolSize int) {
	b.Run(fmt.Sprintf("populationSize=%d poolSize=%d", populationSize, poolSize), func(b *testing.B) {
		fitnesses := makeRandomFitnesses(populationSize)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			randomTournamentWinner(poolSize, fitnesses)
		}
	})
}

func BenchmarkRandomTournamentWinner(b *testing.B) {
	benchRandomTournamentWinner(b, 50, 2)

	benchRandomTournamentWinner(b, 100, 2)
	benchRandomTournamentWinner(b, 1000, 2)

	benchRandomTournamentWinner(b, 50, 5)
	benchRandomTournamentWinner(b, 50, 10)
}
