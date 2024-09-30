package composer

import (
	"io"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Composable[X any] interface {
	Compose(r io.ReadSeeker) (*FileTrie[X], error)
}

type Composer[
	Dic codec.WordDict[T],
	Enc codec.Encoder[X],
	X any,
	T comparable,
] struct {
	newEncoder func(Dic) Enc
	reader     func(io.ReadSeeker, item.Factory, Dic) error
	dict       Dic
}

func NewComposer[
	Dic codec.WordDict[T],
	Enc codec.Encoder[X],
	X any,
	T comparable,
](
	dict Dic,
	newEncoder func(Dic) Enc,
	reader func(io.ReadSeeker, item.Factory, Dic) error,
) *Composer[Dic, Enc, X, T] {
	return &Composer[Dic, Enc, X, T]{
		newEncoder: newEncoder,
		reader:     reader,
		dict:       dict,
	}
}

func (c *Composer[Dic, Enc, X, T]) Compose(
	r io.ReadSeeker,
) (*FileTrie[X], error) {
	f := doublearray.NewBuilder().Factory()
	if err := c.reader(r, f, c.dict); err != nil {
		return nil, err
	}

	da, err := f.Done()
	if err != nil {
		return nil, err
	}
	enc := c.newEncoder(c.dict)

	return NewFileTrie(enc, WithIndex[X](da)), nil
}
