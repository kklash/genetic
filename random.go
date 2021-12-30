package genetic

import (
	"math/rand"
	"sync"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().Unix()))

var randMutex sync.Mutex

func randRangeIntsUnique(max, groupSize int) []int {
	pickedMap := make(map[int]bool)
	for len(pickedMap) < groupSize {
		i := randInt(max)
		pickedMap[i] = true
	}

	picked := make([]int, 0, groupSize)
	for contestantIndex, _ := range pickedMap {
		picked = append(picked, contestantIndex)
	}

	return picked
}

func randInt(max int) int {
	randMutex.Lock()
	n := random.Intn(max)
	randMutex.Unlock()
	return n
}

func randFloat() float64 {
	randMutex.Lock()
	n := random.Float64()
	randMutex.Unlock()
	return n
}
