package genetic_test

import (
	"testing"

	"github.com/kklash/genetic"
)

func pow2(x int) (v int) {
	for v = 1; x > 0; x-- {
		v *= 2
	}
	return v
}

func (problem *KnapsackProblem) PerfectSolution() *KnapsackSolution {
	bestFitness := 0
	bestPackingList := make([]bool, len(problem.Items))
	solution := &KnapsackSolution{
		Problem:     problem,
		PackingList: make([]bool, len(problem.Items)),
	}

	maxIterations := pow2(len(problem.Items))

	for i := 0; i < maxIterations; i++ {
		x := i
		for j := 0; j < len(problem.Items); j++ {
			p := pow2(len(problem.Items) - j - 1)
			if x/p >= 1 {
				x -= p
				solution.PackingList[j] = true
			} else {
				solution.PackingList[j] = false
			}
		}

		fitness := solutionFitness(solution)
		if fitness > bestFitness {
			bestFitness = fitness
			copy(bestPackingList, solution.PackingList)
		}
	}

	solution.PackingList = bestPackingList
	return solution
}

func (solution *KnapsackSolution) Equals(otherSolution *KnapsackSolution) bool {
	if solution.Problem != otherSolution.Problem {
		return false
	}
	if len(solution.PackingList) != len(otherSolution.PackingList) {
		return false
	}
	for i := 0; i < len(solution.PackingList); i++ {
		if solution.PackingList[i] != otherSolution.PackingList[i] {
			return false
		}
	}
	return true
}

func TestAgainstPerfectSolutionAlgorithm(t *testing.T) {
	t.Parallel()

	perfectFitness := 122

	problem := &KnapsackProblem{
		WeightLimit: 31,
		Items: []*KnapsackItem{
			// DO NOT CHANGE without recalculating perfectFitness.
			{Weight: 10, Value: 30},
			{Weight: 5, Value: 20},
			{Weight: 15, Value: 50},
			{Weight: 20, Value: 60},
			{Weight: 2, Value: 10},
			{Weight: 8, Value: 14},
			{Weight: 12, Value: 31},
			{Weight: 14, Value: 45},
			{Weight: 18, Value: 20},
			{Weight: 21, Value: 67},
			{Weight: 1, Value: 7},
			{Weight: 3, Value: 14},
			{Weight: 9, Value: 30},
			{Weight: 19, Value: 45},
			{Weight: 29, Value: 84},
			{Weight: 14, Value: 41},
			{Weight: 22, Value: 35},
			{Weight: 41, Value: 0},
			{Weight: 0, Value: 10},
			{Weight: 31, Value: 103},
		},
	}

	perfectSolutionDone := make(chan struct{})
	evolvedSolutionDone := make(chan struct{})

	go func() {
		perfectSolution := problem.PerfectSolution()
		fitness := solutionFitness(perfectSolution)
		if fitness != perfectFitness {
			t.Errorf("Failed to calculate correct perfectFitness\nWanted %d\nGot    %d", perfectFitness, fitness)
		}
		close(perfectSolutionDone)
	}()

	go func() {
		maxGenerations := 1000
		elitism := 1

		population := genetic.NewPopulation(
			50,
			problem.RandomSolution,
			solutionCrossover,
			genetic.StaticFitnessFunc(solutionFitness),
			genetic.TournamentSelection[*KnapsackSolution](3),
			solutionMutation(0.02),
		)

		expectedAccuracy := 0.9

		minimumFitness := int(float64(perfectFitness)*expectedAccuracy) + 1
		population.Evolve(minimumFitness, maxGenerations, elitism)
		_, evolvedFitness := population.Best()
		accuracy := float64(evolvedFitness) / float64(perfectFitness)
		if accuracy < expectedAccuracy {
			t.Errorf(
				"expected to solve knapsack problem faster than PerfectSolution with accuracy of at least %.2f%%; got %.2f%%",
				expectedAccuracy*100,
				accuracy*100,
			)
		}

		close(evolvedSolutionDone)
	}()

	select {
	case <-perfectSolutionDone:
		t.Errorf("expected evolved solution to be ready before PerfectSolution algorithm.")
		<-evolvedSolutionDone
	case <-evolvedSolutionDone:
		<-perfectSolutionDone
	}
}
