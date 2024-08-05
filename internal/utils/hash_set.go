package utils

import "golang.org/x/exp/maps"

// Definition //////////////////////////////////////////////////////////////////

type HashSet[T comparable] struct {
	table map[T]interface{}
}

// Construction ////////////////////////////////////////////////////////////////

// NewHashSet creates a new set containing any provided elements.
func NewHashSet[T comparable](items ...T) *HashSet[T] {
	return newHashSet[T](len(items)).Add(items...)
}

// NewHashSetFrom creates a new set from provided elements processed by the
// given transform function.
func NewHashSetFrom[T comparable, O any](inputs []O, transform func(input O) (T, error)) (*HashSet[T], error) {
	set := newHashSet[T](len(inputs))
	for _, input := range inputs {
		item, err := transform(input)
		if err != nil {
			return nil, err
		}
		set.put(item)
	}
	return set, nil
}

func NewHashSetAfter[T comparable](callback func(*HashSet[T])) *HashSet[T] {
	set := newHashSet[T](0)
	callback(set)
	return set
}

func NewHashSetAfterErrable[T comparable](callback func(*HashSet[T]) error) (*HashSet[T], error) {
	set := newHashSet[T](0)
	return set, callback(set)
}

func NewHashSetAfterOkable[T comparable](callback func(*HashSet[T]) bool) (*HashSet[T], bool) {
	set := newHashSet[T](0)
	return set, callback(set)
}

func newHashSet[T comparable](expectedItems int) *HashSet[T] {
	if expectedItems < 1 {
		return &HashSet[T]{
			table: make(map[T]interface{}),
		}
	} else {
		return &HashSet[T]{
			table: make(map[T]interface{}, expectedItems),
		}
	}
}

// Inspection //////////////////////////////////////////////////////////////////

// Has returns true if the provided value exists in the current set.
func (set *HashSet[T]) Has(item T) bool {
	_, ok := set.table[item]
	return ok
}

// Size returns the number of elements in the set.
func (set *HashSet[T]) Size() int {
	return len(set.table)
}

// Mutation ////////////////////////////////////////////////////////////////////

// Add adds any provided values to the set and returns itself.
func (set *HashSet[T]) Add(items ...T) *HashSet[T] {
	for _, item := range items {
		set.table[item] = struct{}{}
	}
	return set
}

// Remove removes the provided values from the set and returns itself.
func (set *HashSet[T]) Remove(items ...T) *HashSet[T] {
	for _, item := range items {
		delete(set.table, item)
	}
	return set
}

// Clear empties the set of all values and returns itself.
func (set *HashSet[T]) Clear() *HashSet[T] {
	set.table = make(map[T]interface{})
	return set
}

func (set *HashSet[T]) put(item T) {
	set.table[item] = struct{}{}
}

// Transformation //////////////////////////////////////////////////////////////

// Filter returns a subset, that contains only the values that satisfies the given filter.
func (set *HashSet[T]) Filter(filter func(item T) bool) *HashSet[T] {
	result := NewHashSet[T]()
	for item, _ := range set.table {
		if filter(item) {
			result.table[item] = struct{}{}
		}
	}
	return result
}

// Union creates a new set with all unique values present in either set.
func (set *HashSet[T]) Union(other *HashSet[T]) *HashSet[T] {
	result := NewHashSet[T]()

	// First add all the items from the first table
	for item, _ := range set.table {
		result.table[item] = struct{}{}
	}

	// Then do the same for the second.
	for item, _ := range other.table {
		result.table[item] = struct{}{}
	}

	return result
}

// Intersect creates a new set containing only values that are present in both
// sets.
func (set *HashSet[T]) Intersect(other *HashSet[T]) *HashSet[T] {
	return set.Filter(other.Has)
}

// Without creates a new set with only values from the first set that don't
// appear in the second set.
func (set *HashSet[T]) Without(other *HashSet[T]) *HashSet[T] {
	return set.Filter(func(item T) bool {
		return !other.Has(item)
	})
}

// Conversion //////////////////////////////////////////////////////////////////

// ToSlice returns a slice containing all elements currently in the set in
// indeterminate order.
func (set *HashSet[T]) ToSlice() []T {
	return maps.Keys(set.table)
}
