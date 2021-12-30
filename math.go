package genetic

import (
	"math"
)

func zeroValue[T any]() T {
	var v T
	return v
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func floatDiv(x, y int) float64 {
	if y == 0 {
		return math.Inf(1)
	}

	return float64(x) / float64(y)
}

// computeProportions calculates the proportion each number
// represents in the sum of a given set of integers.
func computeProportions(ints []int) []float64 {
	sum := 0
	for _, n := range ints {
		if n < 0 {
			panic("failed to compute proportion with negative number in set")
		}
		sum += n
	}

	proportions := make([]float64, len(ints))
	for i, n := range ints {
		proportions[i] = floatDiv(n, sum)
	}

	return proportions
}
