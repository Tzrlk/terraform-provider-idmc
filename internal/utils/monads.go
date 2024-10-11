package utils

import . "github.com/samber/mo"

// MapResultOk transforms the value of a successful result to another type.
func MapResultOk[T, R any](from Result[T], fn func(data T) Result[R]) Result[R] {

	// Pass-through if already errored.
	if from.IsError() {
		return Err[R](from.Error())
	}

	// Return the result of the transform.
	return fn(from.MustGet())

}

func MapResultErr[T any](
	from Result[T],
	fn func(err error) Result[T],
) Result[T] {

	// Pass-through if not an error.
	if from.IsOk() {
		return from
	}

	// Otherwise transform the current error.
	return fn(from.Error())

}
