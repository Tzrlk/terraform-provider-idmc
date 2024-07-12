package utils

type HashSet[T any] struct {
	table  map[int]T
	hasher func(T) int
}

func (set *HashSet[T]) Has(item T) bool {
	_, ok := set.table[set.hasher(item)]
	return ok
}

func (set *HashSet[T]) Add(item T) {
	set.table[set.hasher(item)] = item
}

func (set *HashSet[T]) Remove(item T) {
	delete(set.table, set.hasher(item))
}

func (set *HashSet[T]) Clear() {
	set.table = make(map[int]T)
}

func (set *HashSet[T]) Size() int {
	return len(set.table)
}

func NewHashSet[T any](hasher func(T) int) *HashSet[T] {
	return &HashSet[T]{
		table:  make(map[int]T),
		hasher: hasher,
	}
}

//optional functionalities

// AddAll Add multiple values in the set
func (set *HashSet[T]) AddAll(items ...T) {
	for _, item := range items {
		set.Add(item)
	}
}

// Filter returns a subset, that contains only the values that satisfies the given predicate P
func (set *HashSet[T]) Filter(filter func(item T) bool) *HashSet[T] {
	result := NewHashSet(set.hasher)
	for hash, value := range set.table {
		if filter(value) {
			result.table[hash] = value
		}
	}
	return result
}

func (set *HashSet[T]) Union(other *HashSet[T]) *HashSet[T] {
	result := NewHashSet(set.hasher)

	// First add all the items from the first table
	for hash, value := range set.table {
		result.table[hash] = value
	}

	// Then do the same for the second.
	for hash, value := range other.table {
		result.table[hash] = value
	}

	return result
}

func (set *HashSet[T]) Intersect(s2 *HashSet[T]) *HashSet[T] {
	return set.Filter(s2.Has)
}

// Difference returns the subset from s, that doesn't exists in s2 (param)
func (set *HashSet[T]) Difference(s2 *HashSet[T]) *HashSet[T] {
	return set.Filter(func(item T) bool {
		return !s2.Has(item)
	})
}
