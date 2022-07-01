package core

func Values[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
