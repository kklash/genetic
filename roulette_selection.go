package genetic

func rouletteSpin(proportions []float64) int {
	winnerPosition := randFloat()

	for i, proportion := range proportions {
		winnerPosition -= proportion
		if winnerPosition <= 0 {
			return i
		}
	}

	panic("failed to find roulette spin winner")
}

// RouletteSelection spins a virtual roulette wheel to pick mating pairs.
// Higher fitness genomes get proportionally larger sections of the roulette wheel.
// Genomes cannot mate with themselves.
func RouletteSelection[T any](genomes []T, fitnesses []int) [][2]T {
	populationSize := len(genomes)
	flooredFitnesses := make([]int, populationSize)

	for i, fitness := range fitnesses {
		if fitness < 0 {
			fitness = 0
		}
		flooredFitnesses[i] = fitness
	}

	wheelProportions := computeProportions(flooredFitnesses)

	matingPairs := make([][2]T, 0, populationSize/2)
	for len(matingPairs)*2 < populationSize {
		mate1Index := rouletteSpin(wheelProportions)
		mate1 := genomes[mate1Index]

		var mate2Index int
		for mate2Index != mate1Index {
			mate2Index = rouletteSpin(wheelProportions)
		}
		mate2 := genomes[mate2Index]

		matingPairs = append(matingPairs, [2]T{mate1, mate2})
	}

	return matingPairs
}
