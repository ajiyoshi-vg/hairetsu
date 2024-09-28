package dict

import (
	"io"
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Builder[T comparable, X any, D codec.FillableDict[T]] struct {
	items      func(io.Reader) iter.Seq[X]
	units      func(io.Reader) iter.Seq[T]
	newEncoder func(D) codec.Encoder[X]
}

func NewBuilder[T comparable, X any, D codec.FillableDict[T]](
	items func(io.Reader) iter.Seq[X],
	keys func(io.Reader) iter.Seq[T],
	newEncoder func(D) codec.Encoder[X],
) *Builder[T, X, D] {
	return &Builder[T, X, D]{
		items:      items,
		units:      keys,
		newEncoder: newEncoder,
	}
}

func (b *Builder[T, X, D]) Build(r io.ReadSeeker, f item.Factory, dict D) error {
	c := NewCounter(dict)
	c.Add(b.units(r))
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
