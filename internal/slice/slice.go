package slice

func Map[T1, T2 any](s []T1, mp func(T1, int) T2) []T2 {
	r := make([]T2, len(s))

	for i, e := range s {
		r[i] = mp(e, i)
	}

	return r
}

func ForEach[T any](s []T, f func(T, int)) {
	for i, e := range s {
		f(e, i)
	}
}
