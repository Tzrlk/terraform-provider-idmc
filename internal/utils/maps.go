package utils

func TransformMapValues[K comparable, VI any, VO any](input map[K]VI, transform func(K, VI) VO) map[K]VO {
	result := make(map[K]VO)
	for key, val := range input {
		result[key] = transform(key, val)
	}
	return result
}

// MapMerge merges some number of maps together, with subsequent ones
// overwriting entries declared in previous ones. No modifications are made to
// any of the input maps.
func MapMerge[K comparable, V any](items ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, item := range items {
		for k, v := range item {
			result[k] = v
		}
	}
	return result
}
