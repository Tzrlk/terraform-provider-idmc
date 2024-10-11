package utils

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
