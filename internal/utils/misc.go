package utils

func Ptr[T any](val T) *T {
	return &val
}

func Val[T any](ptr *T) T {
	return *ptr
}

func ValOr[T any](ptr *T, or T) T {
	if ptr != nil {
		return *ptr
	}
	return or
}

func OkVal[T any](val T) (T, error) {
	return val, nil
}

func OkPtr[T any](ptr *T) (*T, error) {
	return ptr, nil
}

func TransformMapValues[K comparable, VI any, VO any](input map[K]VI, transform func(K, VI) VO) map[K]VO {
	result := make(map[K]VO)
	for key, val := range input {
		result[key] = transform(key, val)
	}
	return result
}

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

func Handle0(handler func(error), target func() error) func() {
	return func() { handler(target()) }
}
func Handle1[A any](handler func(error), target func(A) error) func(A) {
	return func(a A) { handler(target(a)) }
}
func Handle2[A any, B any](handler func(error), target func(A, B) error) func(A, B) {
	return func(a A, b B) { handler(target(a, b)) }
}
func Handle3[A any, B any, C any](handler func(error), target func(A, B, C) error) func(A, B, C) {
	return func(a A, b B, c C) { handler(target(a, b, c)) }
}
