package genetic_test

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/kklash/genetic"
)

func testCrossoverFunc(t *testing.T, crossover genetic.CrossoverFunc[[]byte]) error {
	male := make([]byte, 32)
	female := make([]byte, 32)

	if _, err := rand.Read(male); err != nil {
		return fmt.Errorf("failed to populate testing male DNA")
	}

	if _, err := rand.Read(female); err != nil {
		return fmt.Errorf("failed to populate testing female DNA")
	}

	child1, child2 := crossover(male, female)

	if len(child1) != 32 || len(child2) != 32 {
		return fmt.Errorf("children DNA have unexpected size")
	}

	for i := 0; i < 32; i++ {
		if child1[i] == male[i] {
			if child2[i] != female[i] {
				return fmt.Errorf("single point crossover failed")
			}
		} else if child1[i] == female[i] {
			if child2[i] != male[i] {
				return fmt.Errorf("single point crossover failed")
			}
		} else {
			return fmt.Errorf("single point crossover failed")
		}
	}
	return nil
}

func TestCrossover(t *testing.T) {
	if err := testCrossoverFunc(t, genetic.NPointCrossover[[]byte](1)); err != nil {
		t.Errorf(err.Error())
	}
	if err := testCrossoverFunc(t, genetic.NPointCrossover[[]byte](2)); err != nil {
		t.Errorf(err.Error())
	}
	if err := testCrossoverFunc(t, genetic.NPointCrossover[[]byte](3)); err != nil {
		t.Errorf(err.Error())
	}
	if err := testCrossoverFunc(t, genetic.NPointCrossover[[]byte](4)); err != nil {
		t.Errorf(err.Error())
	}
	if err := testCrossoverFunc(t, genetic.NPointCrossover[[]byte](5)); err != nil {
		t.Errorf(err.Error())
	}
	if err := testCrossoverFunc(t, genetic.NPointCrossover[[]byte](10)); err != nil {
		t.Errorf(err.Error())
	}

	if err := testCrossoverFunc(t, genetic.UniformCrossover[[]byte]); err != nil {
		t.Errorf(err.Error())
	}

	if err := testCrossoverFunc(t, genetic.AsexualCrossover[[]byte]); err != nil {
		t.Errorf(err.Error())
	}
}
