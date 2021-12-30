package genetic_test

import (
	"crypto/rand"
	"testing"

	"github.com/kklash/bits"
	"github.com/kklash/genetic"
)

type KnapsackItem struct {
	Weight, Value int
}

type KnapsackProblem struct {
	Items       []*KnapsackItem
	WeightLimit int
}

func (problem *KnapsackProblem) RandomSolution() *KnapsackSolution {
	randomBits := make(bits.Bits, len(problem.Items))
	randomBits.ReadFrom(rand.Reader)

	return &KnapsackSolution{
		PackingList: randomBits.Bools(),
		Problem:     problem,
	}
}

type KnapsackSolution struct {
	PackingList []bool
	Problem     *KnapsackProblem
}

func (solution *KnapsackSolution) String() string {
	return bits.BoolsToBits(solution.PackingList).String()
}

func solutionFitness(solution *KnapsackSolution) int {
	if len(solution.PackingList) != len(solution.Problem.Items) {
		panic("received solution with incorrect packing list length")
	}

	weight := 0
	value := 0
	for i, item := range solution.Problem.Items {
		if solution.PackingList[i] {
			weight += item.Weight
			if weight > solution.Problem.WeightLimit {
				return -1
			}
			value += item.Value
		}
	}

	return value
}

func solutionCrossover(s1, s2 *KnapsackSolution) (o1, o2 *KnapsackSolution) {
	o1 = new(KnapsackSolution)
	o2 = new(KnapsackSolution)

	o1.PackingList, o2.PackingList = genetic.SinglePointCrossover(s1.PackingList, s2.PackingList)

	o1.Problem = s1.Problem
	o2.Problem = s1.Problem
	return
}

func solutionMutation(mutationRate float64) genetic.MutationFunc[*KnapsackSolution] {
	mutatePackingList := genetic.RandomizedBinaryMutation(mutationRate)
	return func(solution *KnapsackSolution) {
		mutatePackingList(solution.PackingList)
	}
}

func TestGenetic(t *testing.T) {
	t.Parallel()

	for _, perfectSolution := range knapsackSolutionFixtures {
		perfectFitness := solutionFitness(perfectSolution)
		maxGenerations := 1000
		elitism := 2

		population := genetic.NewPopulation(
			60,
			perfectSolution.Problem.RandomSolution,
			solutionCrossover,
			solutionFitness,
			genetic.TournamentSelection[*KnapsackSolution](3),
			solutionMutation(0.1),
		)

		expectedAccuracy := 0.99

		minimumFitness := int(float64(perfectFitness)*expectedAccuracy) + 1
		population.Evolve(minimumFitness, maxGenerations, elitism)
		bestSolution, bestFitness := population.Best()
		accuracy := float64(bestFitness) / float64(perfectFitness)

		if accuracy < expectedAccuracy {
			t.Errorf("expected to solve knapsack problem with accuracy of at least %.2f%%; got %.2f%%", expectedAccuracy*100, accuracy*100)
			t.Errorf("evolved fitness: %d", bestFitness)
			t.Errorf("evolved solution: %s", bestSolution)
			t.Errorf("perfect fitness: %d", perfectFitness)
			t.Errorf("perfect solution: %s", perfectSolution)
		}
	}
}
