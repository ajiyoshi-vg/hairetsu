package dict

import (
	"io"
	"iter"
	"slices"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Builder[T comparable, X any, D codec.FillableDict[T]] struct {
	items      func(io.Reader) iter.Seq[X]
	units      func(iter.Seq[X]) iter.Seq[T]
	newEncoder func(D) codec.Encoder[X]
}

func NewBuilder[T comparable, X any, D codec.FillableDict[T]](
	items func(io.Reader) iter.Seq[X],
	units func(iter.Seq[X]) iter.Seq[T],
	newEncoder func(D) codec.Encoder[X],
) *Builder[T, X, D] {
	return &Builder[T, X, D]{
		items:      items,
		units:      units,
		newEncoder: newEncoder,
	}
}

func (b *Builder[T, X, D]) Build(r io.ReadSeeker, f item.Factory, dict D) error {
	c := NewCounter(dict)
	c.Add(b.units(b.items(r)))
	c.Build()

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}

	enc := b.newEncoder(dict)
	var i uint32
	for x := range b.items(r) {
		f.Put(item.New(enc.Encode(x), i))
		i++
	}
	return nil
}

func (b *Builder[T, X, D]) BuildSlice(xs []X, f item.Factory, dict D) error {
	c := NewCounter(dict)
	c.Add(b.units(slices.Values(xs)))
	c.Build()

	enc := b.newEncoder(dict)
	for i, x := range xs {
		f.Put(item.New(enc.Encode(x), uint32(i)))
	}
	return nil
}
