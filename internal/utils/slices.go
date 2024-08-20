package utils

func JoinSlices[T any](left []T, right []T) []T {
	result := make([]T, len(left)+len(right))
	_ = copy(result, left)
	_ = copy(result[len(left):], right)
	return result
}

func NewSliceFrom[T any](slice ...[]T) []T {

	// Calculate total required capacity
	var length = 0
	for _, slice := range slice {
		length += len(slice)
	}

	// Create the resulting array with said capacity
	result := make([]T, length)

	// Copy data from each slice into the array.
	var index = 0
	for _, slice := range slice {
		copy(result[index:], slice)
		index += len(slice)
	}

	return result
}

// Coalesce returns the first non-nil pointer in the inputs, or nil.
func Coalesce[T any](items ...*T) *T {
	for _, item := range items {
		if item != nil {
			return item
		}
	}
	return nil
}

func TransformSlice[F any, T any](from []F, to func(from F) T) []T {
	out := make([]T, len(from))
	for index, item := range from {
		out[index] = to(item)
	}
	return out
}
func TransformSliceErr[F any, T any](from []F, to func(from F) (T, error)) ([]T, error) {
	out := make([]T, len(from))
	for index, item := range from {
		val, err := to(item)
		if err != nil {
			return nil, err
		}
		out[index] = val
	}
	return out, nil
}
