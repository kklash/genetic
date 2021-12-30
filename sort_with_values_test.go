package genetic

import (
	"fmt"
	"testing"
)

func copySlice[T any](v []T) []T {
	m := make([]T, len(v))
	copy(m, v)
	return m
}

func TestSortWithValues(t *testing.T) {
	type Fixture struct {
		order        sortOrder
		items        []string
		values       []int
		sortedItems  []string
		sortedValues []int
	}

	fixtures := []*Fixture{
		{
			order:        sortAscending,
			items:        []string{"first", "middle", "last"},
			values:       []int{1, 2, 3},
			sortedItems:  []string{"first", "middle", "last"},
			sortedValues: []int{1, 2, 3},
		},
		{
			order:        sortAscending,
			items:        []string{"middle", "first", "last"},
			values:       []int{2, 1, 3},
			sortedItems:  []string{"first", "middle", "last"},
			sortedValues: []int{1, 2, 3},
		},
		{
			order:        sortDescending,
			items:        []string{"first", "middle", "last"},
			values:       []int{3, 2, 1},
			sortedItems:  []string{"first", "middle", "last"},
			sortedValues: []int{3, 2, 1},
		},
		{
			order:        sortDescending,
			items:        []string{"last", "middle", "first"},
			values:       []int{1, 2, 3},
			sortedItems:  []string{"first", "middle", "last"},
			sortedValues: []int{3, 2, 1},
		},
		{
			order:        sortDescending,
			items:        []string{"last", "first", "middle"},
			values:       []int{1, 3, 2},
			sortedItems:  []string{"first", "middle", "last"},
			sortedValues: []int{3, 2, 1},
		},
	}

	for _, fixture := range fixtures {
		items, values := copySlice(fixture.items), copySlice(fixture.values)
		sortWithValues(fixture.order, items, values)

		itemsString := fmt.Sprintf("%v", items)
		valuesString := fmt.Sprintf("%v", values)
		expectedItemsString := fmt.Sprintf("%v", fixture.sortedItems)
		expectedValuesString := fmt.Sprintf("%v", fixture.sortedValues)

		if itemsString != expectedItemsString {
			t.Errorf("expected items to be sorted\nWanted %s\nGot    %s", expectedItemsString, itemsString)
		}
		if valuesString != expectedValuesString {
			t.Errorf("expected values to be sorted\nWanted %s\nGot    %s", expectedValuesString, valuesString)
		}
	}
}
