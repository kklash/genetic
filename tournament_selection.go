package genetic

func randomTournamentWinner(poolSize int, fitnesses []int) int {
	contestants := randRangeIntsUnique(len(fitnesses), poolSize)

	bestContestant := contestants[0]
	for _, contestant := range contestants[1:] {
		if fitnesses[contestant] > fitnesses[bestContestant] {
			bestContestant = contestant
		}
	}

	return bestContestant
}

// TournamentSelection returns a SelectionFunc of T which selects mates by
// randomly creating 'tournaments' between poolSize contestants. The fittest
// contestants from each random tournament are selected to mate.
// Genomes cannot mate with themselves.
func TournamentSelection[T any](poolSize int) SelectionFunc[T] {
	if poolSize < 2 {
		panic("cannot use tournament selection with pool size less than 2")
	}

	return func(population []T, fitnesses []int) [][2]T {
		populationSize := len(population)

		if poolSize > populationSize {
			panic("cannot select from tournament pool greater than population size")
		}

		matingPairs := make([][2]T, 0, populationSize/2)

		for len(matingPairs)*2 < populationSize {
			var matingPairIndexes [2]int
			for i := 0; i < len(matingPairIndexes); i++ {
				matingPairIndexes[i] = randomTournamentWinner(poolSize, fitnesses)
			}

			if matingPairIndexes[0] != matingPairIndexes[1] {
				matingPairs = append(matingPairs, [2]T{
					population[matingPairIndexes[0]],
					population[matingPairIndexes[1]],
				})
			}
		}

		return matingPairs
	}
}
