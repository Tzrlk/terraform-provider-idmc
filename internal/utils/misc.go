package utils

func Ptr[T any](val T) *T {
	return &val
}

func Val[T any](ptr *T) T {
	return *ptr
}

func OkVal[T any](val T) (T, error) {
	return val, nil
}

func OkPtr[T any](ptr *T) (*T, error) {
	return ptr, nil
}
