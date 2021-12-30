package genetic_test

import (
	"encoding/json"
	"os"
)

var knapsackSolutionFixtures []*KnapsackSolution

func init() {
	fixtures, err := readKnapsackFixtures()
	if err != nil {
		panic(err)
	}

	knapsackSolutionFixtures = fixtures
}

func readKnapsackFixtures() ([]*KnapsackSolution, error) {
	type knapsackProblemJSON struct {
		Items       [][2]int `json:"items"`
		WeightLimit int      `json:"weight_limit"`
		Solution    []int    `json:"solution"`
	}

	var problemsJSON []*knapsackProblemJSON

	data, err := os.ReadFile("knapsack_problems.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &problemsJSON); err != nil {
		return nil, err
	}

	knapsackSolutionFixtures := make([]*KnapsackSolution, len(problemsJSON))
	for i, problemJSON := range problemsJSON {
		solution := &KnapsackSolution{
			Problem: &KnapsackProblem{
				WeightLimit: problemJSON.WeightLimit,
				Items:       make([]*KnapsackItem, len(problemJSON.Items)),
			},
			PackingList: make([]bool, len(problemJSON.Items)),
		}

		for j := 0; j < len(problemJSON.Items); j++ {
			solution.Problem.Items[j] = &KnapsackItem{
				Weight: problemJSON.Items[j][0],
				Value:  problemJSON.Items[j][1],
			}
			solution.PackingList[j] = problemJSON.Solution[j] == 1
		}

		knapsackSolutionFixtures[i] = solution
	}
	return knapsackSolutionFixtures, nil
}
