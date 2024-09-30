package dict

import (
	"io"
	"iter"
	"slices"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Builder[
	T comparable,
	X any,
	Dic codec.FillableDict[T],
	Enc codec.Encoder[X],
] struct {
	items      func(io.Reader) iter.Seq[X]
	units      func(iter.Seq[X]) iter.Seq[T]
	newEncoder func(Dic) Enc
}

func NewBuilder[
	T comparable,
	X any,
	Dic codec.FillableDict[T],
	Enc codec.Encoder[X],
](
	items func(io.Reader) iter.Seq[X],
	units func(iter.Seq[X]) iter.Seq[T],
	newEncoder func(Dic) Enc,
) *Builder[T, X, Dic, Enc] {
	return &Builder[T, X, Dic, Enc]{
		items:      items,
		units:      units,
		newEncoder: newEncoder,
	}
}

func (b *Builder[T, X, Dic, Enc]) Build(
	r io.ReadSeeker,
	f item.Factory,
	dict Dic,
) error {
	c := NewCounter(dict)
	c.Add(b.units(b.items(r)))
	c.Build()

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}

	enc := b.Encoder(dict)
	var i uint32
	for x := range b.items(r) {
		f.Put(item.New(enc.Encode(x), i))
		i++
	}
	return nil
}

func (b *Builder[T, X, Dic, Enc]) BuildSlice(
	xs []X,
	f item.Factory,
	dict Dic,
) error {
	c := NewCounter(dict)
	c.Add(b.units(slices.Values(xs)))
	c.Build()

	enc := b.Encoder(dict)
	for i, x := range xs {
		f.Put(item.New(enc.Encode(x), uint32(i)))
	}
	return nil
}

func (b *Builder[T, X, Dic, Enc]) Encoder(dict Dic) Enc {
	return b.newEncoder(dict)
}
