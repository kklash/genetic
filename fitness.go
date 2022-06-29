package genetic

// StaticFitnessFunc is a utility which maps a static non-competitive fitness function,
// whose output is not dependent on other competing genomes, into a FitnessFunc[T].
// Use this if your genomes' fitnesses are measured independently of the wider population.
func StaticFitnessFunc[T any](fitness func(T) int) FitnessFunc[T] {
	return func(genomes []T, fitnesses []int) {
		for i, genome := range genomes {
			if fitnesses[i] == 0 {
				// Only calculate fitness for genomes whose fitnesses are unknown.
				fitnesses[i] = fitness(genome)
			}
		}
	}
}
