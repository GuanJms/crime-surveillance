package ptr

func Of[T any](v T) *T {
	return &v
}

func DeferOrZero[T any](p *T) T {
	if p != nil {
		return *p
	}
	var zero T
	return zero
}
