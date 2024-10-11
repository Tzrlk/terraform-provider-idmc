package utils

func BadVal[T any](err error, val T) (T, error) {
	return val, err
}

func OkVal[T any](val T) (T, error) {
	return val, nil
}
