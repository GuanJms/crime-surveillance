package ptr

func Of[T any](v T) *T {
	return &v
}

func Deref[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}
