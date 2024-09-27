package doublebyte

import "math"

type Option[T any] func(*T)

func WithContent(m MapDict) Option[ArrayDict] {
	return func(a *ArrayDict) {
		for n := range math.MaxUint16 {
			a[n] = m.Code(uint16(n))
		}
	}
}
