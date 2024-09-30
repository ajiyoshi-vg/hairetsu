package trie

import "github.com/ajiyoshi-vg/hairetsu/doublearray"

type Option[T any] func(*T)

func WithIndex[X any](da *doublearray.DoubleArray) Option[File[X]] {
	return func(dest *File[X]) {
		dest.da = da
	}
}
