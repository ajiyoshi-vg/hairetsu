package codec

import (
	"iter"
)

type Counter[T comparable, S Fillable[T]] struct {
	count map[T]int
	dest  S
}

func NewCounter[T comparable, S Fillable[T]](dest S) *Counter[T, S] {
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

func instantCount[T comparable, S Fillable[T]](dest S, seq iter.Seq[T]) S {
	b := NewCounter(dest)
	b.Add(seq)
	return b.Build()
}
