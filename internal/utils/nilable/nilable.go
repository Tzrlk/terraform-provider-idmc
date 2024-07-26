package nilable

func Or[L any](left *L, right *L) *L {
	if left != nil {
		return left
	}
	return right
}

func OrFn[L any](left *L, right func() *L) *L {
	if left != nil {
		return left
	}
	return right()
}

func OrErr[L any](left *L, right func() error) (*L, error) {
	if left != nil {
		return left, nil
	}
	return nil, right()
}

func OrFnErr[L any](left *L, right func() (*L, error)) (*L, error) {
	if left != nil {
		return left, nil
	}
	return right()
}

func Next[L any, R any](left *L, right func(val *L) *R) *R {
	if left == nil {
		return nil
	}
	return right(left)
}
