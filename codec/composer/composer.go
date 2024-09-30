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
	ComposeFromSlice(xs []X) (*trie.File[X], error)
}

type Composer[
	X any,
	T comparable,
	Dic codec.WordDict[T],
	Enc codec.Encoder[X],
] struct {
	dict    Dic
	builder *dict.Builder[T, X, Dic, Enc]
	option  []doublearray.Option
}

func NewComposer[
	X any,
	T comparable,
	Dic codec.WordDict[T],
	Enc codec.Encoder[X],
](
	dict Dic,
	builder *dict.Builder[T, X, Dic, Enc],
	option ...doublearray.Option,
) *Composer[X, T, Dic, Enc] {
	ret := &Composer[X, T, Dic, Enc]{
		dict:    dict,
		builder: builder,
		option:  option,
	}
	return ret
}

func (c *Composer[X, T, Dic, Enc]) Compose(r io.ReadSeeker) (*trie.File[X], error) {
	f := doublearray.NewBuilder(c.option...).Factory()
	if err := c.builder.Build(r, f, c.dict); err != nil {
		return nil, err
	}

	da, err := f.Done()
	if err != nil {
		return nil, err
	}
	enc := c.builder.Encoder(c.dict)

	return trie.NewFile(enc, trie.WithIndex[X](da)), nil
}

func (c *Composer[X, T, Dic, Enc]) ComposeFromSlice(xs []X) (*trie.File[X], error) {
	f := doublearray.NewBuilder().Factory()
	if err := c.builder.BuildSlice(xs, f, c.dict); err != nil {
		return nil, err
	}

	da, err := f.Done()
	if err != nil {
		return nil, err
	}
	enc := c.builder.Encoder(c.dict)

	return trie.NewFile(enc, trie.WithIndex[X](da)), nil
}
