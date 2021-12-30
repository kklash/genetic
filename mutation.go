package genetic

// RandomizedBinaryMutation returns a MutationFunc which acts on binary genomes (slices of booleans).
// It randomly flips alleles in the mutant genome at the given mutationRate, which should be
// between 0.0 (no mutation) and 1.0 (every allele is flipped).
func RandomizedBinaryMutation(mutationRate float64) MutationFunc[[]bool] {
	if mutationRate <= 0 || mutationRate >= 1 {
		panic("invalid mutation rate, must be between 0 - 1")
	}

	return func(genome []bool) {
		for i := 0; i < len(genome); i++ {
			if randFloat() < mutationRate {
				genome[i] = !genome[i]
			}
		}
	}
}
