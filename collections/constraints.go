package collections

type Ordered interface {
	Integer | Float | ~string
}

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	Signed | Unsigned
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Slice[T any] interface {
	~[]T
}

type Map[K comparable, V any] interface {
	~map[K]V
}

type Chan[T any] interface {
	~chan T
}
