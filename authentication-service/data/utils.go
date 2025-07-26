package data

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func DeferOrZero[T any](p *T) T {
	if p != nil {
		return *p
	}
	var zero T
	return zero
}
