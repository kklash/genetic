package genetic

import (
	"constraints"
	"sort"
)

type sortOrder int

const (
	sortAscending sortOrder = iota
	sortDescending
)

type valuedSlice[T any, V constraints.Ordered] struct {
	Order  sortOrder
	Items  []T
	Values []V
}

// Len implements sort.Interface. Returns the size of the valuedSlice.
func (vs *valuedSlice[T, V]) Len() int {
	return len(vs.Items)
}

// Less implements sort.Interface. Returns true if the fitness of genome i is greater than genome j.
func (vs *valuedSlice[T, V]) Less(i, j int) bool {
	if vs.Order == sortDescending {
		return vs.Values[i] > vs.Values[j]
	} else if vs.Order == sortAscending {
		return vs.Values[i] < vs.Values[j]
	}
	return false
}

// Swap implements sort.Interface. Swaps items of indexes i and j.
func (vs *valuedSlice[T, V]) Swap(i, j int) {
	vs.Items[i], vs.Items[j] = vs.Items[j], vs.Items[i]
	vs.Values[i], vs.Values[j] = vs.Values[j], vs.Values[i]
}

// sortWithValues sorts two slices. The elements of both slices are sorted in ascending or
// descending order, as specified by order, by comparing the elements of the values slice.
func sortWithValues[T any, V constraints.Ordered](order sortOrder, items []T, values []V) {
	sort.Sort(&valuedSlice[T, V]{
		Order:  order,
		Items:  items,
		Values: values,
	})
}
