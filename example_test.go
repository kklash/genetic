package genetic_test

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kklash/genetic"
)

func ExamplePopulation() {
	rand.Seed(time.Now().Unix())

	// We'll try to guess this string using a genetic algorithm.
	bytesToGuess := []byte("how did you ever guess my secret!")

	population := genetic.NewPopulation(
		// population size
		100,

		// GenesisFunc[[]byte] - generates random solutions to initialize the population.
		func() []byte {
			genome := make([]byte, len(bytesToGuess))
			rand.Read(genome)
			return genome
		},

		// CrossoverFunc[[]byte] - combines two parent genomes into two new child genomes.
		genetic.UniformCrossover[[]byte],

		// FitnessFunc[[]byte] - determines how accurate the guess genome is.
		func(guess []byte) (fitness int) {
			fitness = 1
			for i, b := range guess {
				if b == bytesToGuess[i] {
					fitness++
				}
			}
			return
		},

		// SelectionFunc[[]byte] - selects which genomes will reproduce.
		genetic.TournamentSelection[[]byte](3),

		// MutationFunc[[]byte] - randomly alters a given genome to introduce extra variety.
		func(guess []byte) {
			r := make([]byte, len(guess))
			rand.Read(r)
			for i := range guess {
				if rand.Float64() < 0.05 { // 5% mutation rate
					guess[i] ^= r[i]
				}
			}
		},
	)

	population.Evolve(
		len(bytesToGuess)+1, // minimum fitness. Evolution terminates upon finding a solution with this fitness or higher.
		2000,                // max generations. Evolution terminates after iterating through this many generations.
		2,                   // elitism. Carries over this many of the highest-fitness genomes from each previous generation into the next.
	)
	bestSolution, bestFitness := population.Best()

	fmt.Printf("evolved string: %q\n", bestSolution)
	fmt.Println("best fitness:", bestFitness)

	// output:
	// evolved string: "how did you ever guess my secret!"
	// best fitness: 34
}
