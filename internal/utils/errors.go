package utils

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
