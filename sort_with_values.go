package genetic

import (
	"sort"
)

type sortOrder int

const (
	sortAscending sortOrder = iota
	sortDescending
)

type valuedSlice[T any] struct {
	Order  sortOrder
	Items  []T
	Values []int
}

// Len implements sort.Interface. Returns the size of the valuedSlice.
func (vs *valuedSlice[T]) Len() int {
	return len(vs.Items)
}

// Less implements sort.Interface. Returns true if the fitness of genome i is greater than genome j.
func (vs *valuedSlice[T]) Less(i, j int) bool {
	if vs.Order == sortDescending {
		return vs.Values[i] > vs.Values[j]
	} else if vs.Order == sortAscending {
		return vs.Values[i] < vs.Values[j]
	}
	return false
}

// Swap implements sort.Interface. Swaps items of indexes i and j.
func (vs *valuedSlice[T]) Swap(i, j int) {
	vs.Items[i], vs.Items[j] = vs.Items[j], vs.Items[i]
	vs.Values[i], vs.Values[j] = vs.Values[j], vs.Values[i]
}

// sortWithValues sorts two slices. The elements of both slices are sorted in ascending or
// descending order, as specified by order, by comparing the elements of the values slice.
func sortWithValues[T any](order sortOrder, items []T, values []int) {
	sort.Sort(&valuedSlice[T]{
		Order:  order,
		Items:  items,
		Values: values,
	})
}
