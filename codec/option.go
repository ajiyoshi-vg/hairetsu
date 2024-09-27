package codec

type Option[T any] func(*T)

func WithContent[T tinyInteger](m MapDict[T]) Option[ArrayDict[T]] {
	return func(x *ArrayDict[T]) {
		for n := range x.bufferLength() {
			(*x)[n] = m.Code(T(n))
		}
	}
}
