// Package genetic provides genetic algorithms to approximate solutions
// for problems with very large search spaces.
package genetic

import (
	"fmt"
	"reflect"
)

// PopulationSizeMinimum is the minimum size of a Population.
const PopulationSizeMinimum = 2

// GenesisFunc is a function which initializes a randomized possible value for T.
type GenesisFunc[T any] func() T

// CrossoverFunc recombines two genomes to produce two offspring.
// Crossover functions should NOT handle mutation.
type CrossoverFunc[T any] func(male, female T) (T, T)

// FitnessFunc calculates the fitnesses of all genomes in a population,
// storing the results in the given fitnesses slice.
//
// Some entries in fitnesses may be prepopulated - these are cached fitnesses for elite
// genomes surviving from the previous generation. A FitnessFunc may recalculate or
// skip them as needed.
type FitnessFunc[T any] func(allGenomes []T, fitnesses []int)

// MutationFunc randomly alters the DNA of the given genome, in the hopes that
// some mutatations will result in fitter genomes.
type MutationFunc[T any] func(T)

// SelectionFunc selects pairs of mates from a given population of genomes.
// A SelectionFunc is passed a slice of genomes and a corresponding slice of their
// respective pre-computed fitnesses. It should return a slice of mating pairs
// whose length is such that len(matingPairs)*2 >= len(genomes).
//
// A SelectionFunc should NOT mutate the values passed to it.
type SelectionFunc[T any] func(genomes []T, fitnesses []int) (matingPairs [][2]T)

// Population is a struct representing a population of individuals (genomes of
// type T) which can be evolved using genetic algorithms.
type Population[T any] struct {
	genomes   []T
	fitnesses []int

	// Crossover is used to recombine two genomes of type T.
	Crossover CrossoverFunc[T]

	// Fitness computes the fitnesses of a population of genomes of type T.
	Fitness FitnessFunc[T]

	// Selection selects which genomes will reproduce, and which genomes they will mate with.
	Selection SelectionFunc[T]

	// Mutation randomly mutates a genome.
	Mutation MutationFunc[T]
}

// NewPopulation initializes a Population of genomes of the given size.
// The generate function is used to create a genome population of the given size.
func NewPopulation[T any](
	size int,
	generate GenesisFunc[T],
	crossover CrossoverFunc[T],
	fitness FitnessFunc[T],
	selection SelectionFunc[T],
	mutation MutationFunc[T],
) *Population[T] {

	if size < PopulationSizeMinimum {
		panic(fmt.Sprintf("Population size minimum is %d; got %d", PopulationSizeMinimum, size))
	} else if generate == nil {
		panic("expected to receive GenesisFunc")
	} else if crossover == nil {
		panic("expected to receive CrossoverFunc")
	} else if fitness == nil {
		panic("expected to receive FitnessFunc")
	} else if selection == nil {
		panic("expected to receive SelectionFunc")
	}

	population := &Population[T]{
		genomes:   make([]T, size),
		fitnesses: make([]int, size),
		Crossover: crossover,
		Fitness:   fitness,
		Selection: selection,
		Mutation:  mutation,
	}

	for i := 0; i < size; i++ {
		genome := generate()
		population.genomes[i] = genome
	}

	fitness(population.genomes, population.fitnesses)
	sortWithValues(sortDescending, population.genomes, population.fitnesses)

	return population
}

// EvolveOnce evolves the population by one generation, replacing the current population
// with their children. It calls the population's selection function once, its fitness
// function once, and crossover once for every mating pair needed to repopulate.
func (population *Population[T]) EvolveOnce(elitism int) {
	elitism = max(elitism, 0)
	matingPairs := population.Selection(population.genomes, population.fitnesses)

	childGenomes := make([]T, 0, len(matingPairs)*2+elitism)
	if cap(childGenomes) < len(population.genomes) {
		panic("too few mating pairs returned by population's SelectionFunc")
	}

	for _, matingPair := range matingPairs {
		offspring1, offspring2 := population.Crossover(matingPair[0], matingPair[1])
		if population.Mutation != nil {
			population.Mutation(offspring1)
			population.Mutation(offspring2)
		}
		childGenomes = append(childGenomes, offspring1, offspring2)
	}

	nextGenomes := make([]T, len(childGenomes)+elitism)
	copy(nextGenomes, population.genomes[:elitism])
	copy(nextGenomes[elitism:], childGenomes)

	nextFitnesses := make([]int, len(childGenomes)+elitism)
	copy(nextFitnesses, population.fitnesses[:elitism])

	population.Fitness(nextGenomes, nextFitnesses)

	sortWithValues(sortDescending, nextGenomes, nextFitnesses)

	population.genomes = nextGenomes[:len(population.genomes)]
	population.fitnesses = nextFitnesses[:len(population.fitnesses)]
}

// Evolve evolves the population until either a genome is produced which meets the
// given fitnessThreshold, or the maxGenerations threshold is reached.
func (population *Population[T]) Evolve(fitnessThreshold, maxGenerations, elitism int) {
	for i := 0; i < maxGenerations; i++ {
		_, bestFitness := population.Best()
		// fmt.Printf("Gen %d best fitness: %d - diversity %f\n", i, bestFitness, population.Diversity())
		if bestFitness >= fitnessThreshold {
			break
		}

		population.EvolveOnce(elitism)
	}
}

// Best returns the current population's fittest genome and fitness.
func (population *Population[T]) Best() (T, int) {
	return population.genomes[0], population.fitnesses[0]
}

// Diversity compares every genome in the population with one another using reflect.DeepEqual to determine
// whether they are genetic-identicals. It returns a float in range [0.0, 1.0] indicating the percentage of comparisons
// which were NOT identical.
//
// The total number of DeepEqual comparison calls made will be ((s-1)^2 + (s-1)) / 2, where s is the population size.
func (population *Population[T]) Diversity() float64 {
	sames := float64(0)
	opportunities := float64(0)

	for i, g1 := range population.genomes {
		for j := i + 1; j < len(population.genomes); j++ {
			g2 := population.genomes[j]
			opportunities += 1.0
			if reflect.DeepEqual(g1, g2) {
				sames += 1.0
			}
		}
	}

	return 1.0 - sames/opportunities
}
