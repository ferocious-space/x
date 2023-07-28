package collections

func Min[T Ordered](a ...T) (m T) {
	if len(a) > 0 {
		m = a[0]
	}
	for i := 0; i < len(a); i++ {
		if a[i] < m {
			m = a[i]
		}
	}
	return
}

func Max[T Ordered](a ...T) (m T) {
	if len(a) > 0 {
		m = a[0]
	}
	for i := 0; i < len(a); i++ {
		if a[i] > m {
			m = a[i]
		}
	}
	return
}

func ChunkBy[T Slice[N], N any](items T, chunk int) (chunks []T) {
	for chunk < len(items) {
		chunks = append(chunks, items[0:chunk])
		items = items[chunk:]
	}
	return append(chunks, items)
}

func Contains[T comparable](s []T, e T) bool {
	for i := range s {
		if s[i] == e {
			return true
		}
	}
	return false
}

func SliceFilter[T any](s []T, filterFN func(T) bool) (output []T) {
	if filterFN == nil {
		filterFN = func(_ T) bool {
			return true
		}
	}
	for i := range s {
		if filterFN(s[i]) {
			output = append(output, s[i])
		}
	}
	return
}

func SliceApply[T any](s []T, applyFN func(T) T) []T {
	for i := range s {
		s[i] = applyFN(s[i])
	}
	return s
}

func Reduce[T, N any](s []T, reduceFN func(N, T) N, initValue N) N {
	out := initValue
	for _, v := range s {
		out = reduceFN(out, v)
	}
	return out
}

func MapKeys[K comparable, V any](m map[K]V) (out []K) {
	for k := range m {
		out = append(out, k)
	}
	return
}

func MapValues[K comparable, V any](m map[K]V) (out []V) {
	for _, v := range m {
		out = append(out, v)
	}
	return
}

func UniqueNonEmptyElementsOf[T comparable](s []T) (out []T) {
	var zero T
	unique := make(map[T]bool, len(s))
	for _, elem := range s {
		if elem != zero {
			if !unique[elem] {
				out = append(out, elem)
				unique[elem] = true
			}
		}
	}
	return out
}
