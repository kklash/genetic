package genetic

import (
	"fmt"
	"sort"
)

// SinglePointCrossover crosses over two segments of equal-length DNA by choosing a random
// break point in the range of their lengths. It creates two offspring DNA sequences by
// swapping and splicing together the four resulting subslices of DNA.
//
// SinglePointCrossover is analagous to NPointCrossover(1).
func SinglePointCrossover[T ~[]E, E any](male, female T) (T, T) {
	return NPointCrossover[T](1)(male, female)
}

// NPointCrossover crosses two genomes by choosing pointCount random break points in the domain of
// their lengths. It creates two offspring DNA sequences by swapping and splicing together
// the resulting pointCount + 1 subslices of DNA, alternating between male & female parents.
func NPointCrossover[T ~[]E, E any](pointCount int) CrossoverFunc[T] {
	if pointCount < 1 {
		panic(fmt.Sprintf("Invalid point count for NPointCrossover: %d", pointCount))
	}

	return func(male, female T) (T, T) {
		dnaLength := len(male)
		if len(female) != dnaLength {
			panic("cannot do point-based crossover with mismatching DNA length")
		}

		crossoverPoints := make([]int, pointCount+2)
		crossoverPoints[0] = 0
		crossoverPoints[len(crossoverPoints)-1] = dnaLength
		for i := 0; i < pointCount; i++ {
			crossoverPoints[i+1] = randInt(dnaLength)
		}

		sort.Ints(crossoverPoints)

		offspring1 := make(T, dnaLength)
		offspring2 := make(T, dnaLength)

		for i := 1; i < len(crossoverPoints); i++ {
			from, to := crossoverPoints[i-1], crossoverPoints[i]
			if i%2 == 0 {
				copy(offspring1[from:to], male[from:to])
				copy(offspring2[from:to], female[from:to])
			} else {
				copy(offspring1[from:to], female[from:to])
				copy(offspring2[from:to], male[from:to])
			}
		}

		return offspring1, offspring2
	}
}

// UniformCrossover assembles two child genomes from two parents by randomly picking
// individual genome elements from parents. Randomization is exclusive: if one child genome
// inherits one allele from its male parent, the sister genome is guaranteed to inherit the
// opposite allele from the female parent.
func UniformCrossover[T ~[]E, E any](male, female T) (T, T) {
	dnaLength := len(male)
	if len(female) != dnaLength {
		panic("cannot do point-based crossover with mismatching DNA length")
	}

	offspring1 := make(T, dnaLength)
	offspring2 := make(T, dnaLength)

	for i := 0; i < dnaLength; i++ {
		if randFloat() > 0.5 {
			offspring1[i] = male[i]
			offspring2[i] = female[i]
		} else {
			offspring1[i] = female[i]
			offspring2[i] = male[i]
		}
	}

	return offspring1, offspring2
}

// AsexualCrossover does not cross over the given genomes - it simply returns them
// back to the caller, without cloning them.
func AsexualCrossover[T any](male, female T) (T, T) {
	return male, female
}
