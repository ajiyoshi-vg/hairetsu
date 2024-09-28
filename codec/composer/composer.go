package composer

import (
	"io"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Composer[
	X any,
	T comparable,
	Dic codec.WordDict[T],
	Enc codec.Encoder[X],
] struct {
	newEncoder func(Dic) Enc
	reader     func(io.ReadSeeker, item.Factory, Dic) error
	dict       Dic
}

func NewComposer[
	X any,
	T comparable,
	Dic codec.WordDict[T],
	Enc codec.Encoder[X],
](
	dict Dic,
	newEncoder func(Dic) Enc,
	reader func(io.ReadSeeker, item.Factory, Dic) error,
) *Composer[X, T, Dic, Enc] {
	return &Composer[X, T, Dic, Enc]{
		newEncoder: newEncoder,
		reader:     reader,
		dict:       dict,
	}
}

func (c *Composer[X, T, Dic, Enc]) Compose(
	r io.ReadSeeker,
) (*Trie[X, *doublearray.DoubleArray], error) {
	f := doublearray.NewBuilder().Factory()
	if err := c.reader(r, f, c.dict); err != nil {
		return nil, err
	}

	da, err := f.Done()
	if err != nil {
		return nil, err
	}
	enc := c.newEncoder(c.dict)

	return NewTrie(enc, da), nil
}
