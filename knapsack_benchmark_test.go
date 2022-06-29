package genetic_test

import (
	"testing"

	"github.com/kklash/genetic"
)

func benchEvolveOnce(b *testing.B, selection genetic.SelectionFunc[*KnapsackSolution]) {
	problem := knapsackSolutionFixtures[7].Problem // hardest problem
	populationSize := 60
	elitism := 2

	population := genetic.NewPopulation(
		populationSize,
		problem.RandomSolution,
		solutionCrossover,
		genetic.StaticFitnessFunc(solutionFitness),
		selection,
		solutionMutation(0.02),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		population.EvolveOnce(elitism)
	}
}

func BenchmarkPopulation_EvolveOnce(b *testing.B) {
	b.Run("TournamentSelection/poolsize=2", func(b *testing.B) {
		benchEvolveOnce(b, genetic.TournamentSelection[*KnapsackSolution](2))
	})
	b.Run("TournamentSelection/poolsize=5", func(b *testing.B) {
		benchEvolveOnce(b, genetic.TournamentSelection[*KnapsackSolution](5))
	})
	b.Run("TournamentSelection/poolsize=10", func(b *testing.B) {
		benchEvolveOnce(b, genetic.TournamentSelection[*KnapsackSolution](10))
	})
	b.Run("RouletteSelection", func(b *testing.B) {
		benchEvolveOnce(b, genetic.RouletteSelection[*KnapsackSolution])
	})
}
