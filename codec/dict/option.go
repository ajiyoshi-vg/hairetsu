package dict

type Option[T any] func(*T)

func WithContent[T tinyInteger](m Map[T]) Option[Array[T]] {
	return func(x *Array[T]) {
		for n := range x.bufferLength() {
			(*x)[n] = m.Code(T(n))
		}
	}
}
