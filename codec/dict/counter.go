package dict

import (
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/codec"
)

type Counter[T comparable, S codec.Fillable[T]] struct {
	count map[T]int
	dest  S
}

func NewCounter[T comparable, S codec.Fillable[T]](dest S) *Counter[T, S] {
	return &Counter[T, S]{
		count: make(map[T]int),
		dest:  dest,
	}
}

func (b *Counter[T, S]) Add(seq iter.Seq[T]) {
	for x := range seq {
		b.count[x] += 1
	}
}

func (b *Counter[T, S]) Build() S {
	b.dest.Fill(b.count)
	return b.dest
}

func InstantCount[T comparable, S codec.Fillable[T]](dest S, seq iter.Seq[T]) S {
	b := NewCounter(dest)
	b.Add(seq)
	return b.Build()
}
