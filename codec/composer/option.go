package composer

import "github.com/ajiyoshi-vg/hairetsu/doublearray"

type Option[T any] func(*T)

func WithIndex[X any](da *doublearray.DoubleArray) Option[FileTrie[X]] {
	return func(dest *FileTrie[X]) {
		dest.da = da
	}
}
