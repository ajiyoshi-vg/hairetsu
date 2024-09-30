package composer

import (
	"io"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
	"github.com/ajiyoshi-vg/hairetsu/codec/trie"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
)

type Composable[X any] interface {
	Compose(r io.ReadSeeker) (*trie.File[X], error)
}

type Composer[
	Dic codec.WordDict[T],
	Enc codec.Encoder[X],
	X any,
	T comparable,
] struct {
	newEncoder func(Dic) Enc
	dict       Dic
	builder    *dict.Builder[T, X, Dic]
}

func NewComposer[
	Dic codec.WordDict[T],
	Enc codec.Encoder[X],
	X any,
	T comparable,
](
	dict Dic,
	newEncoder func(Dic) Enc,
	builder *dict.Builder[T, X, Dic],
) *Composer[Dic, Enc, X, T] {
	return &Composer[Dic, Enc, X, T]{
		newEncoder: newEncoder,
		dict:       dict,
		builder:    builder,
	}
}

func (c *Composer[Dic, Enc, X, T]) Compose(
	r io.ReadSeeker,
) (*trie.File[X], error) {
	f := doublearray.NewBuilder().Factory()
	if err := c.builder.Build(r, f, c.dict); err != nil {
		return nil, err
	}

	da, err := f.Done()
	if err != nil {
		return nil, err
	}
	enc := c.newEncoder(c.dict)

	return trie.NewFile(enc, trie.WithIndex[X](da)), nil
}
